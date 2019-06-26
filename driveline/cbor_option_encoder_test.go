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
	"time"
)

func TestEncodeQueryOptions(t *testing.T) {
	t.Run("when nil", func(t *testing.T) {
		expected := []byte{cborUndefined}
		actual := make([]byte, 2)
		length := encodeQueryOptions(actual, 0, nil)
		actual = actual[:length]
		if bytes.Compare(expected, actual[:length]) != 0 {
			t.Fail()
		}
	})
	t.Run("with casID", func(t *testing.T) {
		expected := []byte{
			cborArray | 2,
			cborUnsignedInteger | tagReadID,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
		}

		var options QueryOptions
		options.FromRecordID(RecordID{1, 2, 3, 4, 5, 6, 7, 8,})

		actual := make([]byte, 32)
		length := encodeQueryOptions(actual, 0, &options)
		actual = actual[:length]

		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}
	})
}

func TestSizeOfQueryOptions(t *testing.T) {
	t.Run("when nil", func(t *testing.T) {
		if sizeOfQueryOptions(nil) != 1 {
			t.Fail()
		}
	})
	t.Run("when empty", func(t *testing.T) {
		var options QueryOptions
		if sizeOfQueryOptions(&options) != 1 {
			t.Fail()
		}
	})
	t.Run("w/ recordID", func(t *testing.T) {
		t.Run("from head of stream", func(t *testing.T) {
			var options QueryOptions
			options.FromStreamHead()
			if sz := sizeOfQueryOptions(&options); sz != 11 {
				t.Fail()
			}
		})
		t.Run("from tail of stream", func(t *testing.T) {
			var options QueryOptions
			options.fromStreamTail()
			if sz := sizeOfQueryOptions(&options); sz != 11 {
				t.Fail()
			}
		})
		t.Run("from a specific record", func(t *testing.T) {
			var options QueryOptions
			options.FromRecordID(testRecordID)
			if sz := sizeOfQueryOptions(&options); sz != 11 {
				t.Fail()
			}
		})
	})
}

func TestEncodeStoreOptions(t *testing.T) {
	t.Run("when nil", func(t *testing.T) {
		expected := []byte{cborUndefined}
		actual := make([]byte, 1)
		length := encodeStoreOptions(actual, 0, nil)
		actual = actual[:length]
		if bytes.Compare(expected, actual[:length]) != 0 {
			t.Fail()
		}
	})
	t.Run("when empty", func(t *testing.T) {
		expected := []byte{cborUndefined}
		actual := make([]byte, 1)
		var opts StoreOptions
		length := encodeStoreOptions(actual, 0, &opts)
		actual = actual[:length]
		if bytes.Compare(expected, actual[:length]) != 0 {
			t.Fail()
		}
	})
	t.Run("w/ ttl & CAS", func(t *testing.T) {
		expected := []byte{
			cborArray | 4,
			cborUnsignedInteger | tagStoreCASID,
			cborByteString | 8, 1, 2, 3, 4, 5, 6, 7, 8,
			cborUnsignedInteger | tagStoreTTL,
			cborUnsignedInteger | 25, 7, 208,
		}
		var options StoreOptions
		options.CompareAndSwap(RecordID([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
		options.WithTTL(2 * time.Second)
		actual := make([]byte, 32)
		length := encodeStoreOptions(actual, 0, &options)
		actual = actual[:length]
		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}

	})
}

func TestSizeOfStoreOptions(t *testing.T) {
	t.Run("when nil", func(t *testing.T) {
		var options StoreOptions
		if sizeOfStoreOptions(&options) != 1 {
			t.Fail()
		}
	})
	t.Run("w/ ttl & CAS", func(t *testing.T) {
		var options StoreOptions
		options.CompareAndSwap(RecordID([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
		options.WithTTL(2 * time.Second)
		if sz := sizeOfStoreOptions(&options); sz != 15 {
			t.Fail()
		}
	})
}
