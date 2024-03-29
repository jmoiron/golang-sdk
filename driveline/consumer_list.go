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

type listConsumer struct {
	baseConsumer
	isStream bool
	pattern  string
	handler  func(string)
}

var _ consumer = (*listConsumer)(nil)

func newListConsumer(client *Client, consumerID uint64, isStream bool, pattern string, handler func(string)) *listConsumer {
	return &listConsumer{
		baseConsumer: newBaseConsumer(client, consumerID),
		isStream:     isStream,
		pattern:      pattern,
		handler:      handler,
	}
}

func (c *listConsumer) run(ctx context.Context) error {
	return c.Client.list(c)
}

func (c *listConsumer) onRecords(records []Record) {
	if len(records) != 1 {
		c.onFailure(ErrInvalidServerMessage)
		return
	}
	encoded := records[0].Record
	code := encoded[0]
	if !isArray(code) {
		c.onFailure(ErrInvalidServerMessage)
		return
	}
	if lenCode(code) == 0 {
		c.quit()
		return
	}
	entryCount, encoded, err := decodeNumber(encoded)
	if err != nil {
		c.onFailure(err)
		return
	}
	var entry string
	for i := uint64(0); i < entryCount; i++ {
		entry, encoded, err = decodeString(encoded)
		if err != nil {
			c.onFailure(err)
			return
		}
		c.handler(entry)
	}
}

func (c *listConsumer) onDisconnect() {
	c.onFailure(ErrClosed)
}
