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
	"testing"
)

func TestConsumerRegistration(t *testing.T) {
	t.Run("register/unregister", func(t *testing.T) {
		client, _ := testClient()
		c := newSyncConsumer(client, 1533)
		client.registerConsumer(c)
		if _, exists := client.consumers[c.consumerID()]; !exists {
			t.Fail()
		}
		consumerExists := client.unregisterConsumer(c)
		if _, exists := client.consumers[c.consumerID()]; exists {
			t.Fail()
		}
		if !consumerExists {
			t.Fail()
		}
	})

	t.Run("unregister a non existing consumer", func(t *testing.T) {
		client, _ := testClient()
		c := newSyncConsumer(client, 1533)
		consumerExists := client.unregisterConsumer(c)
		if _, exists := client.consumers[c.consumerID()]; exists {
			t.Fail()
		}
		if consumerExists {
			t.Fail()
		}
	})

}

func TestClient_Append(t *testing.T) {
	client, fws := testClient()
	var actual []byte
	fws.WriteHandler = func(buf []byte) (i int, e error) {
		actual = buf
		return len(buf), nil
	}
	if err := client.Append("test-stream", testRecord); err != nil {
		t.Fail()
	}
	expected := encodeAppendByName("test-stream", testRecord)
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}


