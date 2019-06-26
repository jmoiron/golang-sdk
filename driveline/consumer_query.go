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

var _ consumer = (*queryConsumer)(nil)

type queryConsumer struct {
	baseConsumer
	dql          string
	isContinuous bool
	options      *QueryOptions
	handler      func(*Record)
}

func newQueryConsumer(client *Client, consumerID uint64, dql string, isContinuous bool, options *QueryOptions, handler func(*Record)) *queryConsumer {
	if options == nil {
		options = new(QueryOptions)
	}
	return &queryConsumer{
		baseConsumer: newBaseConsumer(client, consumerID),
		dql:          dql,
		options:      options,
		handler:      handler,
		isContinuous: isContinuous,
	}
}

func (c *queryConsumer) run(ctx context.Context) error {
	return c.Client.query(c)
}

func (c *queryConsumer) onRecords(records []Record) {
	cnt := len(records)
	if !c.isContinuous && cnt == 0 {
		c.quit()
		return
	}
	for i := 0; i < cnt; i++ {
		c.options.fromRecordID = records[i].RecordID
		c.handler(&records[i])
	}
}

func (c *queryConsumer) onReconnect() {
	if err := c.Client.query(c); err != nil {
		c.onFailure(err)
	}
}
