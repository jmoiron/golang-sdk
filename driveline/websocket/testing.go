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
	"context"
)

type FakeWebSocket struct {
	opts         webSocketOptions
	WriteHandler func(buf []byte) (int, error)
}

func NewFakeWebsocket() *FakeWebSocket {
	ws := &FakeWebSocket{
		WriteHandler: func(buf [] byte) (int, error) { return len(buf), nil },
	}
	return ws
}

func (ws *FakeWebSocket) Provide(ctx context.Context, endpoint string, options ...Option) (WebSocket, error) {
	ws.opts.configure(options)
	return ws, nil
}

func (ws *FakeWebSocket) Write(buf []byte) (int, error) {
	return ws.WriteHandler(buf)
}

func (ws *FakeWebSocket) Receive(data []byte) {
	ws.opts.messageHandler(ws, data)
}

func (ws *FakeWebSocket) Disconnect() {
	ws.opts.disconnectHandler(ws)
}

func (ws *FakeWebSocket) Reconnect() {
	ws.opts.connectHandler(ws)
}

func (ws *FakeWebSocket) Fail(err error) {
	ws.opts.failureHandler(ws, err)
}

func (ws *FakeWebSocket) Close() error {
	ws.opts.disconnectHandler(ws)
	ws.opts.failureHandler(ws, ErrConnClosed)
	return nil
}
