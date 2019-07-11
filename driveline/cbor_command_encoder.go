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
	"encoding/binary"
)

func encodeAppendByID(streamID uint64, rec []byte) []byte {
	buf := make([]byte, 5+sizeOfNumber(streamID)+1+sizeOfBytes(rec))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 'a'
	buf[3] = 'p'
	buf[4] = 'p'
	// streamID
	off := encodeNumberWithType(buf, 5, streamID, cborUnsignedInteger)
	// Options
	buf[off] = cborUndefined
	off++
	// Data
	encodeBytesWithType(buf, off, rec, cborByteString)
	// Done
	return buf
}

func encodeAppendByName(stream string, rec []byte) []byte {
	buf := make([]byte, 5+sizeOfBytes([]byte(stream))+1+sizeOfBytes(rec))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 'a'
	buf[3] = 'p'
	buf[4] = 'p'
	// StreamName
	off := encodeBytesWithType(buf, 5, []byte(stream), cborTextString)
	// Options
	buf[off] = cborUndefined
	off++
	// Data
	encodeBytesWithType(buf, off, rec, cborByteString)
	// Done
	return buf
}

func encodeCancel(consumerID uint64) []byte {
	buf := make([]byte, 5+sizeOfNumber(consumerID)+1)
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 'c'
	buf[3] = 'a'
	buf[4] = 'n'
	// ConsumerID
	off := encodeNumberWithType(buf, 5, consumerID, cborUnsignedInteger)
	// Options
	buf[off] = cborUndefined
	// Done
	return buf
}

func encodeDefine(aliasID uint8, streamName string) []byte {
	buf := make([]byte, 5+sizeOfNumber(uint64(aliasID))+sizeOfBytes([]byte(streamName)))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 'd'
	buf[3] = 'e'
	buf[4] = 'f'
	// AliasID
	off := encodeNumberWithType(buf, 5, uint64(aliasID), cborUnsignedInteger)
	// Stream Name
	encodeBytesWithType(buf, off, []byte(streamName), cborTextString)
	// Done
	return buf
}

func encodeList(isStream bool, consumerID uint64, pattern string) []byte {
	buf := make([]byte, 5+sizeOfNumber(consumerID)+1+sizeOfBytes([]byte(pattern)))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 3
	if isStream {
		buf[2] = 's'
		buf[3] = 'l'
		buf[4] = 's'
	} else {
		buf[2] = 'l'
		buf[3] = 's'
		buf[4] = 't'
	}
	// ConsumerID
	off := encodeNumberWithType(buf, 5, consumerID, cborUnsignedInteger)
	// Options
	buf[off] = cborUndefined
	off++
	// Pattern
	encodeBytesWithType(buf, off, []byte(pattern), cborTextString)
	// Done
	return buf
}

func encodeLoad(consumerID uint64, key string) []byte {
	buf := make([]byte, 4+sizeOfNumber(consumerID)+1+sizeOfBytes([]byte(key)))
	_ = buf[3] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 2
	buf[2] = 'l'
	buf[3] = 'd'
	// ConsumerID
	off := encodeNumberWithType(buf, 4, consumerID, cborUnsignedInteger)
	// Options
	buf[off] = cborUndefined
	off++
	// Stream
	encodeBytesWithType(buf, off, []byte(key), cborTextString)
	// Done
	return buf
}

func encodeQuery(isContinuous bool, consumerID uint64, dql string, options *QueryOptions) []byte {
	buf := make([]byte, 4+sizeOfNumber(consumerID)+sizeOfQueryOptions(options)+sizeOfBytes([]byte(dql)))
	_ = buf[3] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 2
	if isContinuous {
		buf[2] = 's'
	} else {
		buf[2] = 'q'
	}
	buf[3] = 'q'
	// ConsumerID
	off := encodeNumberWithType(buf, 4, consumerID, cborUnsignedInteger)
	// Options
	off = encodeQueryOptions(buf, off, options)
	// Query
	encodeBytesWithType(buf, off, []byte(dql), cborTextString)
	// Done
	return buf
}

func encodeRemove(key string) []byte {
	buf := make([]byte, 4+1+sizeOfBytes([]byte(key)))
	_ = buf[3] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 2
	buf[2] = 'r'
	buf[3] = 'm'
	// Options
	buf[4] = cborUndefined
	// Key
	encodeBytesWithType(buf, 5, []byte(key), cborTextString)
	// Done
	return buf
}

func encodeRemoveMatches(pattern string) []byte {
	buf := make([]byte, 5+1+sizeOfBytes([]byte(pattern)))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 'r'
	buf[3] = 'm'
	buf[4] = 'k'
	// Options
	buf[5] = cborUndefined
	// Stream Name
	encodeBytesWithType(buf, 6, []byte(pattern), cborTextString)
	// Done
	return buf
}

func encodeStore(key string, data []byte, options *StoreOptions) []byte {
	buf := make([]byte, 4+sizeOfBytes([]byte(key))+sizeOfStoreOptions(options)+sizeOfBytes(data))
	_ = buf[3] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 4
	// Command
	buf[1] = cborTextString | 2
	buf[2] = 's'
	buf[3] = 't'
	// Key
	off := encodeBytesWithType(buf, 4, []byte(key), cborTextString)
	// Options
	off = encodeStoreOptions(buf, off, options)
	// Data
	encodeBytesWithType(buf, off, data, cborByteString)
	// Done
	return buf
}

func encodeSync(consumerID uint64) []byte {
	buf := make([]byte, 5+sizeOfNumber(consumerID))
	_ = buf[4] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 2
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 's'
	buf[3] = 'y'
	buf[4] = 'n'
	// consumerId
	encodeNumberWithType(buf, 5, consumerID, cborUnsignedInteger)
	// Done
	return buf
}

func encodeTruncateByID(stream uint64) []byte {
	buf := make([]byte, 5+1+sizeOfNumber(stream))
	_ = buf[5] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 't'
	buf[3] = 'r'
	buf[4] = 'c'
	// Options
	buf[5] = cborUndefined
	// Stream
	encodeNumberWithType(buf, 6, stream, cborUnsignedInteger)
	// Done
	return buf
}

func encodeTruncateByName(stream string) []byte {
	buf := make([]byte, 5+1+sizeOfBytes([]byte(stream)))
	_ = buf[5] // bounds check elimination
	// Envelope
	buf[0] = cborArray | 3
	// Command
	buf[1] = cborTextString | 3
	buf[2] = 't'
	buf[3] = 'r'
	buf[4] = 'c'
	// Options
	buf[5] = cborUndefined
	// Stream
	encodeBytesWithType(buf, 6, []byte(stream), cborTextString)
	// Done
	return buf
}

func sizeOfNumber(n uint64) int {
	switch {
	case n < 24:
		return 1
	case n < 0x100:
		return 2
	case n < 0x10000:
		return 3
	case n < 0x100000000:
		return 5
	default:
		return 9
	}
}

func sizeOfBytes(b []byte) int {
	n := len(b)
	switch {
	case n < 24:
		return n + 1
	case n < 0x100:
		return n + 2
	case n < 0x10000:
		return n + 3
	case n < 0x100000000:
		return n + 5
	default:
		return n + 9
	}
}

func encodeNumberWithType(buf []byte, off int, n uint64, cborType byte) int {
	buf = buf[off:]
	switch {
	case n < 24:
		buf[0] = byte(n) | cborType
		return off + 1
	case n < 0x100:
		buf[0] = 24 | cborType
		buf[1] = byte(n)
		return off + 2
	case n < 0x10000:
		buf[0] = 25 | cborType
		binary.BigEndian.PutUint16(buf[1:], uint16(n))
		return off + 3
	case n < 0x100000000:
		buf[0] = 26 | cborType
		binary.BigEndian.PutUint32(buf[1:], uint32(n))
		return off + 5
	default:
		buf[0] = 27 | cborType
		binary.BigEndian.PutUint64(buf[1:], n)
		return off + 9
	}
}

func encodeBytesWithType(buf []byte, off int, data []byte, cborType byte) int {
	b := buf[off:]
	n := len(data)
	switch {
	case n == 0:
		b[0] = cborType
		return off + 1
	case n < 24:
		b[0] = byte(n) | cborType
		off += 1
	case n < 0x100:
		b[0] = 24 | cborType
		b[1] = byte(n)
		off += 2
	case n < 0x10000:
		b[0] = 25 | cborType
		binary.BigEndian.PutUint16(b[1:], uint16(n))
		off += 3
	case n < 0x100000000:
		b[0] = 26 | cborType
		binary.BigEndian.PutUint32(b[1:], uint32(n))
		off += 5
	default:
		b[0] = 27 | cborType
		binary.BigEndian.PutUint64(b[1:], uint64(n))
		off += 9
	}
	return off + copy(buf[off:], data)
}

func encodeRecordID(buf []byte, off int, recordID RecordID) int {
	return encodeBytesWithType(buf, off, []byte(recordID), cborByteString)
}
