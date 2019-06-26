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

package driveline

import (
	"context"
	"fmt"
	"sync"

	ws "github.com/1533-systems/golang-sdk/driveline/websocket"
)

type consumer interface {
	consumerID() uint64
	run(ctx context.Context) error
	onRecords([]Record)
	onReconnect()
	onDisconnect()
	onFailure(err error)

	done() <-chan struct{}
	err() error
}

// Client is the connection object that interacts with Driveline
type Client struct {
	endpoint       string
	consumers      map[uint64]consumer
	consumersLock  sync.Mutex
	consumerID     uint64
	consumerIDLock sync.Mutex
	ws             ws.WebSocket
	defines        defines
	errorHandler   func(error)
}

// NewClient creates a new Client
func NewClient(ctx context.Context, endpoint string, options ...option) (*Client, error) {
	c := &Client{
		endpoint:     endpoint,
		consumerID:   0,
		consumers:    make(map[uint64]consumer),
		errorHandler: func(error) {},
	}
	c.defines.reset()
	opts := clientOptions{
		client:       c,
		newWebSocket: ws.New,
	}
	for _, opt := range options {
		opt(&opts)
	}
	opts.wsOptions = append(opts.wsOptions,
		ws.OnConnect(c.onConnect),
		ws.OnDisconnect(c.onDisconnect),
		ws.OnFailure(c.onFailure),
		ws.OnMessage(c.onMessage),
		ws.ErrorHandler(c.errorHandler),
	)
	var err error
	if c.ws, err = opts.newWebSocket(ctx, c.endpoint, opts.wsOptions...); err != nil {
		return nil, err
	}
	return c, nil
}

// Close the connection to Driveline
func (c *Client) Close() error {
	for _, consumer := range c.snapshotConsumers() {
		consumer.onFailure(ErrClosed)
	}
	if c.ws != nil {
		return c.ws.Close()
	}
	return nil
}

// OpenStream creates a Stream object that help save data bandwidth when
// dealing with streams that have a large number of small messages.
func (c *Client) OpenStream(name string) (*Stream, error) {
	id := c.defines.allocate(name)
	if id.isNumeric() {
		if err := c.define(name, uint8(id.numericID())); err != nil {
			return nil, err
		}
	}
	s := &Stream{
		client:   c,
		streamID: id,
	}
	return s, nil
}

// CloseStream closes a Stream object and releases its resources.
func (c *Client) CloseStream(stream *Stream) {
	if stream.streamID != nil {
		c.defines.release(stream.streamID)
		stream.streamID = nil
	}
	if stream.client != nil {
		stream.client = nil
	}
}

func (c *Client) nextConsumerID() uint64 {
	c.consumerIDLock.Lock()
	id := c.consumerID
	c.consumerID += 1
	c.consumerIDLock.Unlock()
	return id
}

func (c *Client) runConsumer(ctx context.Context, consumer consumer) error {
	c.registerConsumer(consumer)
	defer c.unregisterConsumer(consumer)
	if err := consumer.run(ctx); err != nil {
		return err
	}
	select {
	case <-consumer.done():
		return consumer.err()
	case <-ctx.Done():
		if err := c.cancel(consumer); err != nil {
			c.errorHandler(err)
		}
		return ctx.Err()
	}
}

func (c *Client) registerConsumer(consumer consumer) {
	c.consumersLock.Lock()
	c.consumers[consumer.consumerID()] = consumer
	c.consumersLock.Unlock()
}

func (c *Client) unregisterConsumer(consumer consumer) bool {
	c.consumersLock.Lock()
	consumerID := consumer.consumerID()
	_, exists := c.consumers[consumerID]
	if exists {
		delete(c.consumers, consumerID)
	}
	c.consumersLock.Unlock()
	return exists
}

func (c *Client) sendMessage(message []byte) error {
	if _, err := c.ws.Write(message); err != nil {
		return fmt.Errorf("cannot send message: %s", err.Error())
	}
	return nil
}

func (c *Client) onMessage(buf []byte) {
	msg, err := decodeServerMessage(buf)
	if err != nil {
		c.errorHandler(fmt.Errorf("cannot decode server message: %s", err.Error()))
		return
	}
	c.consumersLock.Lock()
	consumer, ok := c.consumers[msg.consumerID]
	c.consumersLock.Unlock()
	if !ok {
		c.errorHandler(fmt.Errorf("received message for unknown consumer %d", msg.consumerID))
		return
	}
	if msg.err != nil {
		consumer.onFailure(msg.err)
		return
	}
	consumer.onRecords(msg.records)
}

func (c *Client) onConnect() {
	for id, stream := range c.defines.aliasMap {
		if err := c.define(stream, id); err != nil {
			c.errorHandler(fmt.Errorf("cannot set alias %d for stream %s", id, stream))
		}
	}
	for _, consumer := range c.snapshotConsumers() {
		consumer.onReconnect()
	}
}

func (c *Client) onDisconnect() {
	for _, consumer := range c.snapshotConsumers() {
		consumer.onDisconnect()
	}
}

func (c *Client) onFailure(err error) {
	for _, consumer := range c.snapshotConsumers() {
		consumer.onFailure(err)
	}
}

func (c *Client) snapshotConsumers() []consumer {
	c.consumersLock.Lock()
	consumers := make([]consumer, len(c.consumers))
	i := 0
	for _, c := range c.consumers {
		consumers[i] = c
		i++
	}
	c.consumersLock.Unlock()
	return consumers
}
