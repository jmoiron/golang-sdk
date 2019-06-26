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

type loadConsumer struct {
	baseConsumer
	key    string
	record *Record
}

var _ consumer = (*loadConsumer)(nil)

func newLoadConsumer(client *Client, consumerID uint64, key string) *loadConsumer {
	return &loadConsumer{
		baseConsumer: newBaseConsumer(client, consumerID),
		key:          key,
	}
}

func (c *loadConsumer) run(ctx context.Context) error {
	return c.Client.load(c)
}

func (c *loadConsumer) onRecords(records []Record) {
	if len(records) != 1 {
		c.onFailure(ErrInvalidServerMessage)
		return
	}
	c.record = &records[0]
	c.quit()
}

func (c *loadConsumer) onReconnect() {
	if err := c.Client.load(c); err != nil {
		c.onFailure(err)
	}
}
