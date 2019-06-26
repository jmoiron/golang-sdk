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
	"math/bits"
)

var recordIDLen = sizeOfBytes([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

func sizeOfQueryOptions(options *QueryOptions) int {
	if options == nil || options.assigned == 0 {
		return 1
	}
	return 1 + 1 + recordIDLen
}

func encodeQueryOptions(buf []byte, off int, options *QueryOptions) int {
	if options == nil || options.assigned == 0 {
		buf[off] = cborUndefined
		return off + 1
	}
	buf[off] = cborArray | 2
	buf[off+1] = encodedReadIDTag
	off = encodeRecordID(buf, off+2, options.fromRecordID)
	return off
}

func sizeOfStoreOptions(options *StoreOptions) int {
	if options == nil || options.assigned == 0 {
		return 1
	}
	size := 1
	if options.assigned&optStoreTTLOption != 0 {
		size += 1 + sizeOfNumber(options.ttl)
	}
	if options.assigned&optStoreCASOption != 0 {
		size += 1 + recordIDLen
	}
	return size
}

func encodeStoreOptions(buf []byte, off int, options *StoreOptions) int {
	if options == nil || options.assigned == 0 {
		buf[off] = cborUndefined
		return off + 1
	}
	buf[off] = cborArray | byte(bits.OnesCount16(uint16(options.assigned))*2)
	off++
	if options.assigned&optStoreCASOption != 0 {
		buf[off] = encodedStoreCASIDTag
		off = encodeRecordID(buf, off+1, options.casRecordID)
	}
	if options.assigned&optStoreTTLOption != 0 {
		buf[off] = encodedStoreTTLTag
		off = encodeNumberWithType(buf, off+1, options.ttl, cborUnsignedInteger)
	}
	return off
}
