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
	"bytes"
	"context"
	"testing"
)

func TestQueryConsumer(t *testing.T) {
	t.Run("sends command on run", func(t *testing.T) {
		client, fws := testClient()
		fws.WriteHandler = func(buf []byte) (int, error) {
			if bytes.Compare(buf, encodeQuery(true, 1533, "pattern", nil)) != 0 {
				t.Fail()
			}
			return len(buf), nil
		}
		c := newQueryConsumer(client, 1533, "pattern", true, nil, func(r *Record) {
			t.Fail()
		})
		if err := c.run(context.Background()); err != nil {
			t.Fail()
		}
	})

	t.Run("receives records", func(t *testing.T) {
		client, _ := testClient()
		var result []*Record
		c := newQueryConsumer(client, 1533, "pattern", true, nil, func(r *Record) {
			result = append(result, r)
		})
		c.onRecords([]Record{
			{Record: []byte{cborTextString | 5, 'h', 'e', 'l', 'l', 'o'}},
			{Record: []byte{cborTextString | 5, 'w', 'o', 'r', 'l', 'd'}},
		})
		if len(result) != 2 {
			t.Fail()
		}
		if bytes.Compare(result[0].Record, []byte{cborTextString | 5, 'h', 'e', 'l', 'l', 'o',}) != 0 {
			t.Fail()
		}
		if bytes.Compare(result[1].Record, []byte{cborTextString | 5, 'w', 'o', 'r', 'l', 'd',}) != 0 {
			t.Fail()
		}
	})

	t.Run("handles disconnections", func(t *testing.T) {
		client, fws := testClient()
		var commandWritten bool
		fws.WriteHandler = func(buf []byte) (int, error) {
			if bytes.Compare(buf, encodeQuery(true, 1533, "pattern", nil)) != 0 {
				t.Fail()
			}
			commandWritten = true
			return len(buf), nil
		}
		c := newQueryConsumer(client, 1533, "pattern", true, nil, func(r *Record) {
			t.Fail()
		})
		c.onDisconnect()
		if c.ctx.Err() != nil {
			t.Fail()
		}
		c.onReconnect()
		if !commandWritten {
			t.Fail()
		}
	})

	t.Run("stops when one-shot query terminates", func(t *testing.T) {
		client, _ := testClient()
		c := newQueryConsumer(client, 1533, "pattern", false, nil, func(r *Record) {
			t.Fail()
		})
		c.onRecords(nil)
		if c.ctx.Err() != context.Canceled {
			t.Fail()
		}

		if c.err() != nil {
			t.Fail()
		}
	})

	t.Run("fails when reconnection fails", func(t *testing.T) {
		client, fws := testClient()
		var commandWritten bool
		fws.WriteHandler = func(buf []byte) (int, error) {
			if bytes.Compare(buf, encodeQuery(true, 1533, "pattern", nil)) != 0 {
				t.Fail()
			}
			commandWritten = true
			return 0, ErrClosed
		}
		c := newQueryConsumer(client, 1533, "pattern", true, nil, func(r *Record) {
			t.Fail()
		})
		c.onDisconnect()
		if c.ctx.Err() != nil {
			t.Fail()
		}
		c.onReconnect()
		if !commandWritten {
			t.Fail()
		}
		if c.err() == nil {
			t.Fail()
		}
	})
}
