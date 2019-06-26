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
	"strings"
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	t.Run("only decode an array", func(t *testing.T) {
		sample := []byte{
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborArray | 2,
			cborUnsignedInteger | tagMessageID,
			cborArray | 1,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
			cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
	t.Run("fails to decode an invalid array", func(t *testing.T) {
		sample := []byte{
			cborArray | 28,
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborArray | 2,
			cborUnsignedInteger | tagMessageID,
			cborArray | 1,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
			cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
	t.Run("fails to decode a unknown command", func(t *testing.T) {
		sample := []byte{
			cborArray | 4,
			cborTextString | 5, 'w', 'r', 'o', 'n', 'g',
			cborUnsignedInteger | 5,
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeDataMessage(t *testing.T) {
	t.Run("decode a valid message w/ RecordID", func(t *testing.T) {
		sample := []byte{
			cborArray | 4,
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborArray | 2,
			cborUnsignedInteger | tagMessageID,
			cborArray | 1,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
			cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
		}
		record, err := decodeServerMessage(sample)
		if err != nil {
			t.Errorf("decoding failed: %s", err)
			t.FailNow()
		}
		if record == nil {
			t.Log("received a nil record")
			t.Fail()
		}
		if record.err != nil {
			t.Log("received a record with an error")
			t.Fail()
		}

		if record.consumerID != 5 {
			t.Fail()
		}

		if len(record.records) != 1 {
			t.Fail()
		}

		if bytes.Compare([]byte{1, 2, 3, 4, 5, 6, 7, 8}, []byte(record.records[0].RecordID)) != 0 {
			t.Fail()
		}
		if bytes.Compare([]byte("payload"), record.records[0].Record) != 0 {
			t.Fail()
		}
	})

	t.Run("decode a valid message w/o RecordID", func(t *testing.T) {
		sample := []byte{
			cborArray | 4,
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborUndefined,
			cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
		}
		record, err := decodeServerMessage(sample)
		if err != nil {
			t.Errorf("decoding failed: %s", err)
			t.FailNow()
		}
		if record == nil {
			t.Log("received a nil record")
			t.Fail()
		}
		if record.err != nil {
			t.Log("received a record with an error")
			t.Fail()
		}

		if record.consumerID != 5 {
			t.Fail()
		}

		if len(record.records) != 1 {
			t.Fail()
		}

		if len(record.records[0].RecordID) != 0 {
			t.Fail()
		}
		if bytes.Compare([]byte("payload"), record.records[0].Record) != 0 {
			t.Fail()
		}
	})

	t.Run("decode a valid message w/o data", func(t *testing.T) {
		sample := []byte{
			cborArray | 4,
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborUndefined,
			cborUndefined,
		}
		record, err := decodeServerMessage(sample)
		if err != nil {
			t.Errorf("decoding failed: %s", err)
			t.FailNow()
		}
		if record == nil {
			t.Log("received a nil record")
			t.Fail()
		}
		if record.err != nil {
			t.Log("received a record with an error")
			t.Fail()
		}

		if record.consumerID != 5 {
			t.Fail()
		}
		if len(record.records) != 0 {
			t.Fail()
		}
	})

	t.Run("fails to decode a malformed message", func(t *testing.T) {

		samples := []struct {
			title  string
			sample []byte
		}{
			{"options invalid", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 4,
				cborUnsignedInteger | tagMessageID,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"consumerID type", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborTextString | 5, 'w', 'r', 'o', 'n', 'g',
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"consumerID encoding", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 28, 0, 0, 0, 0, 0, 0, 0,
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"tag format", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 28, 0, 0, 0, 0, 0, 0, 0,
				cborMap | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"tag array encoding", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 28, 0, 0, 0, 0, 0, 0, 0,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"tag array length", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 3,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"tag code type", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborSignedInteger | tagMessageID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"tag unknown", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborUnsignedInteger | tagReadID,
				cborArray | 1,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"message id type", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"message id encoding", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 28, 9, 9, 9, 9,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"message id number not matching the number of records", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 2,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
			{"invalid record ID", []byte{
				cborArray | 4,
				cborTextString | 4, 'd', 'a', 't', 'a',
				cborUnsignedInteger | 1,
				cborArray | 2,
				cborUnsignedInteger | tagMessageID,
				cborArray | 1,
				cborTextString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
				cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
			},},
		}

		/*

			cborArray | 4,
			cborTextString | 4, 'd', 'a', 't', 'a',
			cborUnsignedInteger | 5,
			cborArray | 2,
			cborUnsignedInteger | tagMessageID,
			cborArray | 1,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
			cborByteString | 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd',
		*/
		for _, s := range samples {
			t.Run(s.title, func(t *testing.T) {
				if _, err := decodeServerMessage(s.sample); err == nil {
					t.Fail()
				}
			})
		}
	})
}

func TestDecodeErrorMessage(t *testing.T) {
	t.Run("decodes a message", func(t *testing.T) {
		sample := []byte{
			cborArray | 3,
			cborTextString | 3, 'e', 'r', 'r',
			cborUnsignedInteger | 5,
			cborTextString | 4, 'o', 'u', 'c', 'h',
		}

		record, err := decodeServerMessage(sample)
		if err != nil {
			t.Log("decoding failed")
			t.Fail()
		}
		if record.err == nil {
			t.Log("received a record without errors")
			t.Fail()
		}

		if record.records != nil {
			t.Fail()
		}
		if strings.Compare("ouch", record.err.Error()) != 0 {
			t.Fail()
		}
	})

	t.Run("fails to decode a malformed consumer id", func(t *testing.T) {
		sample := []byte{
			cborArray | 3,
			cborTextString | 3, 'e', 'r', 'r',
			cborUnsignedInteger | 28,
			cborTextString | 4, 'o', 'u', 'c', 'h',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})

	t.Run("fails to decode a wrong consumer id", func(t *testing.T) {
		sample := []byte{
			cborArray | 3,
			cborTextString | 3, 'e', 'r', 'r',
			cborTextString | 4, 'm', 'y', 'i', 'd',
			cborTextString | 4, 'o', 'u', 'c', 'h',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})

	t.Run("fails to decode the wrong payload type", func(t *testing.T) {
		sample := []byte{
			cborArray | 3,
			cborTextString | 3, 'e', 'r', 'r',
			cborUnsignedInteger | 5,
			cborByteString | 4, 'o', 'u', 'c', 'h',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeSyncMessage(t *testing.T) {
	t.Run("decodes a message", func(t *testing.T) {
		sample := []byte{
			cborArray | 2,
			cborTextString | 3, 's', 'y', 'n',
			cborUnsignedInteger | 5,
		}

		record, err := decodeServerMessage(sample)
		if err != nil {
			t.Fail()
		}
		if record == nil {
			t.Fail()
		}
		if record.consumerID != 5 {
			t.Fail()
		}
		if len(record.records) > 0 {
			t.Fail()
		}
	})

	t.Run("fails to decode a message when the consumerID is invalid", func(t *testing.T) {
		sample := []byte{
			cborArray | 2,
			cborTextString | 3, 's', 'y', 'n',
			cborTextString | 4, 'd', 'a', 't', 'a',
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
	t.Run("fails to decode a message when the consumerID is invalid", func(t *testing.T) {
		sample := []byte{
			cborArray | 2,
			cborTextString | 3, 's', 'y', 'n',
			cborUnsignedInteger | 28, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
		if _, err := decodeServerMessage(sample); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeNumber(t *testing.T) {
	t.Run("decodes a small number", func(t *testing.T) {
		data := []byte{cborUnsignedInteger | 23}
		n, _, err := decodeNumber(data)
		if err != nil {
			t.Fail()
		}
		if n != 23 {
			t.Fail()
		}
	})

	t.Run("decodes an 8-bit  number", func(t *testing.T) {
		data := []byte{cborUnsignedInteger | 24, 200}
		n, _, err := decodeNumber(data)
		if err != nil {
			t.Fail()
		}
		if n != 200 {
			t.Fail()
		}
	})

	t.Run("decodes a 16-bit  number", func(t *testing.T) {
		data := []byte{cborUnsignedInteger | 25, 0x12, 0x34}
		n, _, err := decodeNumber(data)
		if err != nil {
			t.Fail()
		}
		if n != 0x1234 {
			t.Fail()
		}
	})

	t.Run("decodes a 32-bit  number", func(t *testing.T) {
		data := []byte{cborUnsignedInteger | 26, 0x10, 0x20, 0x30, 0x40}
		n, _, err := decodeNumber(data)
		if err != nil {
			t.Fail()
		}
		if n != 0x10203040 {
			t.Fail()
		}
	})
	t.Run("decodes a 64-bit  number", func(t *testing.T) {
		data := []byte{cborUnsignedInteger | 27, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}
		n, _, err := decodeNumber(data)
		if err != nil {
			t.Fail()
		}
		if n != 0x1020304050607080 {
			t.Fail()
		}
	})
}

func TestDecodeString(t *testing.T) {

}
