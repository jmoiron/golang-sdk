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

const (
	cborUnsignedInteger = byte(0 << 5)
	cborSignedInteger   = byte(1 << 5)
	cborByteString      = byte(2 << 5)
	cborTextString      = byte(3 << 5)
	cborArray           = byte(4 << 5)
	cborMap             = byte(5 << 5)
	cborMulti           = byte(7 << 5)

	// cborTag             = byte(6 << 5)

	cborNull      = cborMulti | 22
	cborUndefined = cborMulti | 23

	cborTypeMask   = 0x07 << 5
	cborLengthMask = 0x1f

	tagMessageID  = 1
	tagReadID     = 2
	tagStoreCASID = 3
	tagStoreTTL   = 4

	encodedMessageIdTag  = cborUnsignedInteger | tagMessageID
	encodedReadIDTag     = cborUnsignedInteger | tagReadID
	encodedStoreCASIDTag = cborUnsignedInteger | tagStoreCASID
	encodedStoreTTLTag   = cborUnsignedInteger | tagStoreTTL
)

func lenCode(b byte) uint64 {
	return uint64(b & cborLengthMask)
}

func isArray(b byte) bool {
	return b&cborTypeMask == cborArray
}

func isByteString(b byte) bool {
	return b&cborTypeMask == cborByteString
}

func isTextString(b byte) bool {
	return b&cborTypeMask == cborTextString
}

func isUnsignedInteger(b byte) bool {
	return b&cborTypeMask == cborUnsignedInteger
}

func isUndefined(b byte) bool {
	return b == cborUndefined
}

func isBlank(b byte) bool {
	return b == cborUndefined || b == cborNull
}
