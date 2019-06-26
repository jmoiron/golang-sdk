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
)

type baseConsumer struct {
	Client     *Client
	ConsumerID uint64
	ctx        context.Context
	quit       context.CancelFunc
	result     error
}

func newBaseConsumer(client *Client, consumerID uint64) baseConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	c := baseConsumer{
		Client:     client,
		ConsumerID: consumerID,
		ctx:        ctx,
		quit:       cancel,
	}
	return c
}

func (c *baseConsumer) consumerID() uint64 {
	return c.ConsumerID
}

func (c *baseConsumer) onDisconnect() {
}

func (c *baseConsumer) onReconnect() {
}

func (c *baseConsumer) onFailure(err error) {
	c.result = err
	c.quit()
}

func (c *baseConsumer) done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *baseConsumer) err() error {
	return c.result
}
