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

import "testing"

func BenchmarkDecodeNumber(b *testing.B) {
	b.Run("small number < 24", func(b *testing.B) {
		data := []byte{23}
		b.SetBytes(int64(len(data)))
		for i := 0; i < b.N; i++ {
			decodeNumber(data)
		}
	})
	b.Run("8 bit number", func(b *testing.B) {
		data := []byte{24, 0x01}
		b.SetBytes(int64(len(data)))
		for i := 0; i < b.N; i++ {
			decodeNumber(data)
		}
	})
	b.Run("16 bit number", func(b *testing.B) {
		data := []byte{24, 0x01, 0x02}
		b.SetBytes(int64(len(data)))
		for i := 0; i < b.N; i++ {
			decodeNumber(data)
		}
	})
	b.Run("32 bit number", func(b *testing.B) {
		data := []byte{26, 0x01, 0x02, 0x03, 0x04}
		b.SetBytes(int64(len(data)))
		for i := 0; i < b.N; i++ {
			decodeNumber(data)
		}
	})
	b.Run("64 bit number", func(b *testing.B) {
		data := []byte{27, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
		b.SetBytes(int64(len(data)))
		for i := 0; i < b.N; i++ {
			decodeNumber(data)
		}
	})
}
