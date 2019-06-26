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
	"strconv"
	"testing"
	"time"
)

var payloads = []int{16, 256, 1024, 4096, 65535, 1024 * 1024}

func BenchmarkEncodeAppendByID(b *testing.B) {
	runWithPayload(b, func(b *testing.B, rec []byte) {
		for i := 0; i < b.N; i++ {
			encodeAppendByID(1533, rec)
		}
	})
}

func BenchmarkEncodeAppendByName(b *testing.B) {
	runWithPayload(b, func(b *testing.B, rec []byte) {
		for i := 0; i < b.N; i++ {
			encodeAppendByName("stream-name", rec)
		}
	})
}

func BenchmarkEncodeCancel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeCancel(1533)
	}
}

func BenchmarkEncodeDefine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeDefine(25, "stream-name")
	}
}

func BenchmarkEncodeList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeList(true, 1533, "pattern")
	}
}

func BenchmarkEncodeLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeLoad(1533, "key-name")
	}
}

func BenchmarkEncodeQuery(b *testing.B) {
	b.Run("without-options", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			encodeQuery(true, 1533, "SELECT * FROM 'kv/**'", nil)
		}
	})

	b.Run("with-options", func(b *testing.B) {
		var options QueryOptions
		options.FromStreamHead()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			encodeQuery(true, 1533, "SELECT * FROM 'kv/**'", &options)
		}
	})
}

func BenchmarkEncodeRemove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeRemove("key-name")
	}
}

func BenchmarkEncodeRemoveMatches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeRemoveMatches("key-pattern**")
	}
}

func BenchmarkEncodeStore(b *testing.B) {
	b.Run("without-options", func(b *testing.B) {
		runWithPayload(b, func(b *testing.B, rec []byte) {
			for i := 0; i < b.N; i++ {
				encodeStore("key-name", rec, nil)
			}
		})
	})
	b.Run("with-options", func(b *testing.B) {
		var options StoreOptions
		options.WithTTL(time.Second)
		options.CompareAndSwap(RecordID{1, 2, 3, 4, 5, 6, 7, 8})
		runWithPayload(b, func(b *testing.B, rec []byte) {
			for i := 0; i < b.N; i++ {
				encodeStore("key-name", rec, &options)
			}
		})
	})
}

func BenchmarkEncodeSync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeSync(1533)
	}
}

func BenchmarkTruncateById(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeTruncateByID(1533)
	}
}

func BenchmarkTruncateByName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeTruncateByName("stream-name")
	}
}

func BenchmarkEncodeNumberWithType(b *testing.B) {
	var buf [20]byte
	for i := 0; i < b.N; i++ {
		encodeNumberWithType(buf[0:20], 0, 0x0102030405060708, 0)
	}
}

func runWithPayload(b *testing.B, code func(*testing.B, []byte)) {
	for _, sz := range payloads {
		rec := make([]byte, sz)
		b.Run(strconv.Itoa(sz)+"-byte-payload", func(b *testing.B) {
			b.ResetTimer()
			b.SetBytes(int64(sz))
			code(b, rec)
		})
	}
}
