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
	"reflect"
	"testing"
)

func TestStream_Append(t *testing.T) {
	client, fws := testClient()
	var commands [][]byte
	fws.WriteHandler = func(buf []byte) (int, error) {
		commands = append(commands, buf)
		return len(buf), nil
	}

	stream, err := client.OpenStream("test_stream")
	if err != nil {
		t.Fail()
	}
	if stream == nil {
		t.Fail()
	}
	if len(commands) != 1 {
		t.Fail()
	}
	commands = commands[:0]

	expected := []byte{cborTextString | 4, 'd', 'a', 't', 'a'}
	if err := stream.Append(expected); err != nil {
		t.Fail()
	}
	if len(commands) != 1 {
		t.Fail()
	}
	if !reflect.DeepEqual(commands[0], encodeAppendByID(stream.streamID.numericID(), expected)) {
		t.Fail()
	}
}

func TestStream_Truncate(t *testing.T) {
	client, fws := testClient()
	var commands [][]byte
	fws.WriteHandler = func(buf []byte) (int, error) {
		commands = append(commands, buf)
		return len(buf), nil
	}
	stream, err := client.OpenStream("test_stream")
	if err != nil {
		t.Fail()
	}
	if stream == nil {
		t.Fail()
	}
	if len(commands) != 1 {
		t.Fail()
	}
	commands = commands[:0]
	if err := stream.Truncate(); err != nil {
		t.Fail()
	}
	if len(commands) != 1 {
		t.Fail()
	}
	if !reflect.DeepEqual(commands[0], encodeTruncateByID(stream.streamID.numericID())) {
		t.Fail()
	}
}
