// Copyright 2019, 1533 Systems, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package websocket

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrConnClosed            = errors.New("connection closed")
	ErrInvalidWebSocketFrame = errors.New("invalid WebSocket frame")
	ErrMaxReconnect          = errors.New("maximum reconnection attempts reached")
	ErrUnexpectedEndOfStream = errors.New("unexpected end of stream")
	ErrInvalidFrameType      = errors.New("unexpected frame type received")
)

type frameOpCode byte

const (
	defaultMaxInFlight      = 100
	maxReconnectWaitRatio   = 5
	defaultConnectTimeout   = 2 * time.Second
	defaultReconnectWait    = 1 * time.Second
	defaultMaxReconnectWait = defaultReconnectWait * maxReconnectWaitRatio

	maxOutputBuffer = 32 * 1024 * 1024
	readBufferSize  = 1024*1024 + 65536 + 1024
	maxInputBuffer  = 16 * 1024 * 1024

	continuationFrame = frameOpCode(0x00)
	textFrame         = frameOpCode(0x01)
	binaryFrame       = frameOpCode(0x02)
	closeFrame        = frameOpCode(0x08)
	pingFrame         = frameOpCode(0x09)
	pongFrame         = frameOpCode(0x0A)
)

type WebSocket interface {
	io.WriteCloser
}

var _ WebSocket = (*webSocket)(nil)

type webSocket struct {
	endpoint   string
	outputLock sync.Locker
	cnx        io.ReadWriteCloser
	readBuffer []byte
	closeErr   error
	dataFrames chan []byte
	pongFrames chan []byte
	cancel     func()
	webSocketOptions
}

func New(ctx context.Context, endpoint string, options ...Option) (WebSocket, error) {
	ws := &webSocket{
		endpoint:   endpoint,
		readBuffer: make([]byte, readBufferSize),
		outputLock: new(sync.Mutex),
		pongFrames: make(chan []byte, 1),
		cancel:     func() {},
	}
	ws.webSocketOptions.configure(options)
	ws.dataFrames = make(chan []byte, ws.maxInFlight)
	if err := ws.run(ctx); err != nil {
		ws.Close()
		return nil, err
	}
	return ws, nil
}

func (ws *webSocket) Close() error {
	if ws.isClosed() {
		return nil
	}
	ws.endpoint = ""
	ws.cancel()
	return nil
}

func (ws *webSocket) Write(buf []byte) (int, error) {
	ws.dataFrames <- buf
	return len(buf), nil
}

var errInterrupted = errors.New("goroutine interrupted")

func (ws *webSocket) run(ctx context.Context) error {
	startResult := make(chan error, 1)
	defer close(startResult)
	loopCtx, cancel := context.WithCancel(context.Background())
	ws.cancel = cancel
	go ws.wsLoop(loopCtx, startResult)
	select {
	case <-ctx.Done():
		if err := ws.Close(); err != nil {
			ws.errorHandler(err)
		}
		return ctx.Err()
	case err := <-startResult:
		return err
	}
}

func (ws *webSocket) isClosed() bool {
	return len(ws.endpoint) == 0
}

func (ws *webSocket) wsLoop(ctx context.Context, startResult chan<- error) {
	var (
		err       error
		errWriter error
		errReader error
	)
	for attempt := 0; attempt < ws.maxReconnect || ws.maxReconnect == -1; attempt += 1 {
		select {
		case <-ctx.Done():
			ws.failureHandler(ErrConnClosed)
			return
		case <-time.After(ws.timeWaitForAttempt(attempt)):
			break
		}
		if err = ws.connect(ws.endpoint); err != nil {
			continue
		}
		if startResult != nil {
			startResult <- nil
			startResult = nil
		}
		attempt = 0

		ws.connectHandler()
		runnerCtx, cancelRunner := context.WithCancel(ctx)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer func() {
				if crash := recover(); crash != nil {
					errWriter = fmt.Errorf("recovered writer panic %v", crash)
				}
				cancelRunner()
				wg.Done()
			}()
			errWriter = ws.runWriterLoop(runnerCtx.Done())
		}()
		go func() {
			defer func() {
				if crash := recover(); crash != nil {
					errReader = fmt.Errorf("recovered reader panic %v", crash)
				}
				cancelRunner()
				wg.Done()
			}()
			errReader = ws.runReaderLoop(runnerCtx.Done())
		}()
		wg.Wait()

		if errWriter != nil && errWriter != errInterrupted {
			ws.errorHandler(errWriter)
		}
		if errReader != nil && errReader != errInterrupted {
			ws.errorHandler(errReader)
		}
		ws.disconnectHandler()
		if len(ws.endpoint) == 0 {
			return
		}
	}
	if startResult != nil {
		startResult <- ErrMaxReconnect
		startResult = nil
	}
	ws.failureHandler(ErrMaxReconnect)
}

func (ws *webSocket) timeWaitForAttempt(attempt int) time.Duration {
	timeWait := time.Duration(2^attempt-1) * ws.reconnectWait / 2
	if timeWait > ws.maxReconnectWait {
		timeWait = ws.maxReconnectWait
	}
	return timeWait
}

func (ws *webSocket) runReaderLoop(stopCh <-chan struct{}) error {
	for {
		select {
		case <-stopCh:
			return errInterrupted
		default:
			break
		}
		opCode, frame, err := readFrame(ws.cnx)
		if err != nil {
			return err
		}
		switch opCode {
		case binaryFrame:
			ws.messageHandler(frame)
		case closeFrame:
			// TODO: must send a Close frame
		case pingFrame:
			ws.pongFrames <- frame
		case pongFrame:
			// TODO: add ping/pong support?
		default:
			ws.errorHandler(ErrInvalidFrameType)
		}
	}
}

func (ws *webSocket) runWriterLoop(stopCh <-chan struct{}) error {
	out := bufio.NewWriterSize(ws.cnx, maxOutputBuffer)
	var frame []byte
	for {
		select {
		case <-stopCh:
			return errInterrupted
		case frame = <-ws.pongFrames:
			if err := writeFrame(out, binaryFrame, frame); err != nil {
				return err
			}
		case frame = <-ws.dataFrames:
			if err := writeFrame(out, binaryFrame, frame); err != nil {
				return err
			}
			if messageCnt := len(ws.dataFrames); messageCnt > 0 {
				for i := 0; i < messageCnt; i++ {
					frame = <-ws.dataFrames
					if err := writeFrame(out, binaryFrame, frame); err != nil {
						return err
					}
				}
			}
			if err := out.Flush(); err != nil {
				return err
			}
		}
	}
}

func readFrame(in io.Reader) (frameOpCode, []byte, error) {
	var hdr [14]byte
	read, err := in.Read(hdr[:2])
	if err != nil {
		return 0, nil, err
	}
	if read < 2 {
		return 0, nil, ErrUnexpectedEndOfStream
	}

	opCode := frameOpCode(0x7F & hdr[0])

	fin := hdr[0]&0x80 != 0
	isMasked := (hdr[1] & 0x80) != 0

	if isMasked || !fin {
		return 0, nil, ErrInvalidWebSocketFrame
	}

	frameLen := uint64(0X7F & hdr[1])
	switch frameLen {
	case 126:
		read, err = in.Read(hdr[2:4])
		if err != nil {
			return 0, nil, err
		}
		frameLen = uint64(binary.BigEndian.Uint16(hdr[2:]))

	case 127:
		read, err = in.Read(hdr[2:10])
		if err != nil {
			return 0, nil, err
		}
		frameLen = binary.BigEndian.Uint64(hdr[2:])
	}

	frame := make([]byte, frameLen)
	if _, err := io.ReadFull(in, frame); err != nil {
		return 0, nil, err
	}

	return opCode, frame, nil
}

func writeFrame(out io.Writer, opCode frameOpCode, frame []byte) error {
	var hdr [10]byte
	l := uint64(len(frame))

	hdr[0] = 0x80 | byte(opCode)
	switch {
	case l < 126:
		hdr[1] = byte(l)
		if _, err := out.Write(hdr[0:2]); err != nil {
			return err
		}
	case l < 0x10000:
		hdr[1] = 126
		binary.BigEndian.PutUint16(hdr[2:], uint16(l))
		if _, err := out.Write(hdr[0:4]); err != nil {
			return err
		}
	default:
		hdr[1] = 127
		binary.BigEndian.PutUint64(hdr[2:], l)
		if _, err := out.Write(hdr[0:10]); err != nil {
			return err
		}
	}
	if _, err := out.Write(frame); err != nil {
		return err
	}
	return nil
}
