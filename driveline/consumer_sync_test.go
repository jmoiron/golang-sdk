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

func TestSyncConsumer(t *testing.T) {
	t.Run("sends command on run", func(t *testing.T) {
		client, fws := testClient()
		fws.WriteHandler = func(buf []byte) (int, error) {
			if bytes.Compare(buf, encodeSync(1533)) != 0 {
				t.Fail()
			}
			return len(buf), nil
		}
		c := newSyncConsumer(client, 1533)
		if err := c.run(context.Background()); err != nil {
			t.Fail()
		}
	})

	t.Run("receives a message", func(t *testing.T) {
		client, _ := testClient()
		c := newSyncConsumer(client, 1533)
		c.onRecords([]Record{{Record: nil}})
		if c.err() != nil {
			t.Fail()
		}
		if c.ctx.Err() != context.Canceled {
			t.Fail()
		}
	})

	t.Run("fails when the client disconnects", func(t *testing.T) {
		client, _ := testClient()
		c := newSyncConsumer(client, 1533)
		c.onDisconnect()
		if c.err() == nil {
			t.Fail()
		}
		if c.ctx.Err() == nil {
			t.Fail()
		}
	})

}
