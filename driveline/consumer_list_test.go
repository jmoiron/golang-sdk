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

func TestListConsumer(t *testing.T) {
	t.Run("sends command on run", func(t *testing.T) {
		client, fws := testClient()
		fws.WriteHandler = func(buf []byte) (int, error) {
			if bytes.Compare(buf, encodeList(true, 1533, "pattern")) != 0 {
				t.Fail()
			}
			return len(buf), nil
		}
		c := newListConsumer(client, 1533, true, "pattern", func(s string) {
			t.Fail()
		})
		if err := c.run(context.Background()); err != nil {
			t.Fail()
		}
	})

	t.Run("receives text entries", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{{
			Record: []byte{
				cborArray | 2,
				cborTextString | 5, 'h', 'e', 'l', 'l', 'o',
				cborTextString | 5, 'w', 'o', 'r', 'l', 'd',
			},
		}})
		if len(result) != 2 {
			t.Fail()
		}
		if result[0] != "hello" || result[1] != "world" {
			t.Fail()
		}
	})
	t.Run("stops when it receives a record containing no data", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{{
			Record: []byte{
				cborArray | 0,
			},
		}})
		if c.ctx.Err() != context.Canceled {
			t.Fail()
		}
	})
	t.Run("stops when the client disconnects", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{{
			Record: []byte{
				cborArray | 0,
			},
		}})
		c.onDisconnect()
		if c.ctx.Err() != context.Canceled {
			t.Fail()
		}
	})
	t.Run("fails with a payload containing multiple records", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{
			{
				Record: []byte{
					cborArray | 2,
					cborTextString | 5, 'h', 'e', 'l', 'l', 'o',
					cborTextString | 5, 'w', 'o', 'r', 'l', 'd',
				},
			},
			{
				Record: []byte{
					cborArray | 2,
					cborTextString | 5, 'h', 'e', 'l', 'l', 'o',
					cborTextString | 5, 'w', 'o', 'r', 'l', 'd',
				},
			},
		})
		if len(result) != 0 {
			t.Fail()
		}
		if c.err() == nil {
			t.Fail()
		}
	})
	t.Run("fails with a payload that does not have the right type", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{
			{
				Record: []byte{
					cborTextString | 5, 'h', 'e', 'l', 'l', 'o',
				},
			},
		})
		if len(result) != 0 {
			t.Fail()
		}
		if c.err() == nil {
			t.Fail()
		}
	})
	t.Run("fails with a payload where the number of items is incorrectly encoded", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{
			{
				Record: []byte{
					cborArray | 28,
					cborTextString | 5, 'h', 'e', 'l', 'l', 'o',
					cborTextString | 5, 'w', 'o', 'r', 'l', 'd',
				},
			},
		})
		if len(result) != 0 {
			t.Fail()
		}
		if c.err() == nil {
			t.Fail()
		}
	})
	t.Run("fails with a payload where the items are not strings", func(t *testing.T) {
		client, _ := testClient()
		var result []string
		c := newListConsumer(client, 1533, true, "pattern", func(res string) {
			result = append(result, res)
		})
		c.onRecords([]Record{
			{
				Record: []byte{
					cborArray | 2,
					cborByteString | 5, 'h', 'e', 'l', 'l', 'o',
					cborByteString | 5, 'w', 'o', 'r', 'l', 'd',
				},
			},
		})
		if len(result) != 0 {
			t.Fail()
		}
		if c.err() == nil {
			t.Fail()
		}
	})
}
