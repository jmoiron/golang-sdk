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

	ws "github.com/1533-systems/golang-sdk/driveline/websocket"
)

var testRecordID = RecordID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

var testRecord = []byte{
	cborMap | 1,
	cborTextString | 3, 'k', 'e', 'y',
	cborUnsignedInteger | 5,
}

func testClient() (*Client, *ws.FakeWebSocket) {
	fake := ws.NewFakeWebsocket()
	c, _ := NewClient(context.Background(), "ws://test", websocketProvider(fake.Provide))
	return c, fake
}
