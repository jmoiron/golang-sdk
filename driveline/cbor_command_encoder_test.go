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

func TestEncodeAppendByID(t *testing.T) {
	expected := []byte{
		cborArray | 4,
		cborTextString | 3, 'a', 'p', 'p',
		cborUnsignedInteger | 11,
		cborUndefined,
		cborByteString | 4, 'd', 'a', 't', 'a',
	}
	actual := encodeAppendByID(11, []byte("data"))
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeAppendByName(t *testing.T) {
	expected := []byte{
		cborArray | 4,
		cborTextString | 3, 'a', 'p', 'p',
		cborTextString | 11, 't', 'e', 's', 't', '-', 's', 't', 'r', 'e', 'a', 'm',
		cborUndefined,
		cborByteString | 4, 'd', 'a', 't', 'a',
	}
	actual := encodeAppendByName("test-stream", []byte("data"))
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeDefine(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 3, 'd', 'e', 'f',
		cborUnsignedInteger | 2,
		cborTextString | 11, 't', 'e', 's', 't', '-', 's', 't', 'r', 'e', 'a', 'm',
	}
	actual := encodeDefine(2, "test-stream")
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeQuery(t *testing.T) {
	t.Run("continuous", func(t *testing.T) {
		expected := []byte{
			cborArray | 4,
			cborTextString | 2, 's', 'q',
			cborUnsignedInteger | 24, 123,
			cborUndefined,
			cborTextString | 20, 'S', 'E', 'L', 'E', 'C', 'T', ' ', '*', ' ', 'F', 'R', 'O', 'M', ' ', 's', 't', 'r', 'e', 'a', 'm',
		}
		actual := encodeQuery(true, 123, "SELECT * FROM stream", nil)
		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}
	})

	t.Run("quick", func(t *testing.T) {
		expected := []byte{
			cborArray | 4,
			cborTextString | 2, 'q', 'q',
			cborUnsignedInteger | 24, 123,
			cborUndefined,
			cborTextString | 20, 'S', 'E', 'L', 'E', 'C', 'T', ' ', '*', ' ', 'F', 'R', 'O', 'M', ' ', 's', 't', 'r', 'e', 'a', 'm',
		}
		actual := encodeQuery(false, 123, "SELECT * FROM stream", nil)
		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}

	})
}

func TestEncodeCancel(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 3, 'c', 'a', 'n',
		cborUnsignedInteger | 5,
		cborUndefined,
	}
	actual := encodeCancel(5)
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeTruncateByName(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 3, 't', 'r', 'c',
		cborUndefined,
		cborTextString | 9, 'm', 'y', '-', 's', 't', 'r', 'e', 'a', 'm',
	}
	actual := encodeTruncateByName("my-stream")
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeTruncateByID(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 3, 't', 'r', 'c',
		cborUndefined,
		cborUnsignedInteger | 16,
	}
	actual := encodeTruncateByID(16)
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeRemove(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 2, 'r', 'm',
		cborUndefined,
		cborTextString | 6, 's', 't', 'r', 'e', 'a', 'm',
	}
	actual := encodeRemove("stream")
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeRemoveMatches(t *testing.T) {
	expected := []byte{
		cborArray | 3,
		cborTextString | 3, 'r', 'm', 'k',
		cborUndefined,
		cborTextString | 4, 's', 't', 'r', '*',
	}
	actual := encodeRemoveMatches("str*")
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeLoad(t *testing.T) {
	expected := []byte{
		cborArray | 4,
		cborTextString | 2, 'l', 'd',
		cborUnsignedInteger | 3,
		cborUndefined,
		cborTextString | 6, 's', 't', 'r', 'e', 'a', 'm',
	}
	actual := encodeLoad(3, "stream")
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeStore(t *testing.T) {
	expected := []byte{
		cborArray | 4,
		cborTextString | 2, 's', 't',
		cborTextString | 6, 's', 't', 'r', 'e', 'a', 'm',
		cborUndefined,
		cborByteString | 4, 'd', 'a', 't', 'a',
	}
	actual := encodeStore("stream", []byte("data"), nil)
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestEncodeList(t *testing.T) {
	t.Run("stream", func(t *testing.T) {
		expected := []byte{
			cborArray | 4,
			cborTextString | 3, 's', 'l', 's',
			cborUnsignedInteger | 7,
			cborUndefined,
			cborTextString | 7, 's', 't', 'r', 'e', 'a', 'm', '*',
		}
		actual := encodeList(true, 7, "stream*")
		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}

	})

	t.Run("keys", func(t *testing.T) {
		expected := []byte{
			cborArray | 4,
			cborTextString | 3, 'l', 's', 't',
			cborUnsignedInteger | 7,
			cborUndefined,
			cborTextString | 4, 'k', 'v', '/', '*',
		}
		actual := encodeList(false, 7, "kv/*")
		if bytes.Compare(expected, actual) != 0 {
			t.Fail()
		}
	})
}

func TestEncodeSync(t *testing.T) {
	expected := [] byte{
		cborArray | 2,
		cborTextString | 3, 's', 'y', 'n',
		cborUnsignedInteger | 7,
	}
	actual := encodeSync(7)
	if bytes.Compare(expected, actual) != 0 {
		t.Fail()
	}
}

func TestSizeOfNumber(t *testing.T) {
	t.Run("small number", func(t *testing.T) {
		if sizeOfNumber(22) != 1 {
			t.Fatalf("invalid number size")
		}
	})
	t.Run("8-bit number", func(t *testing.T) {
		if sizeOfNumber(0x07f) != 2 {
			t.Fatalf("invalid number size")
		}

	})

	t.Run("16-bit number", func(t *testing.T) {
		if sizeOfNumber(0x0102) != 3 {
			t.Fatalf("invalid number size")
		}

	})
	t.Run("32-bit number", func(t *testing.T) {
		if sizeOfNumber(0x01020304) != 5 {
			t.Fatalf("invalid number size")
		}

	})
	t.Run("64-bit number", func(t *testing.T) {
		if sizeOfNumber(0x0102030405060708) != 9 {
			t.Fatalf("invalid number size")
		}
	})
}

func TestEncodeNumberWithType(t *testing.T) {
	t.Run("small number", func(t *testing.T) {
		var data = make([]byte, 16)
		off := encodeNumberWithType(data, 0, 23, cborUnsignedInteger)
		if off != 1 {
			t.Fail()
		}
		if data[0] != cborUnsignedInteger|23 {
			t.Fail()
		}
	})
	t.Run("small 8-bit", func(t *testing.T) {
		var data = make([]byte, 16)
		off := encodeNumberWithType(data, 0, 200, cborUnsignedInteger)
		if off != 2 {
			t.Fail()
		}
		if data[0] != cborUnsignedInteger|24 && data[1] != 200 {
			t.Fail()
		}
	})

	t.Run("16-bit", func(t *testing.T) {
		var data = make([]byte, 16)
		off := encodeNumberWithType(data, 0, 0x1234, cborUnsignedInteger)
		if off != 3 {
			t.Fail()
		}
		if bytes.Compare(data[:off], []byte{cborUnsignedInteger | 25, 0x12, 0x34}) != 0 {
			t.Fail()
		}
	})

	t.Run("32-bit", func(t *testing.T) {
		var data = make([]byte, 16)
		off := encodeNumberWithType(data, 0, 0x12345678, cborUnsignedInteger)
		if off != 5 {
			t.Fail()
		}
		if bytes.Compare(data[:off], []byte{cborUnsignedInteger | 26, 0x12, 0x34, 0x56, 0x78}) != 0 {
			t.Fail()
		}
	})

	t.Run("64-bit", func(t *testing.T) {
		var data = make([]byte, 16)
		off := encodeNumberWithType(data, 0, 0x0102030405060708, cborUnsignedInteger)
		if off != 9 {
			t.Fail()
		}
		if bytes.Compare(data[:off], []byte{cborUnsignedInteger | 27, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}) != 0 {
			t.Fail()
		}
	})
}

func TestSizeOfBytes(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		if sizeOfBytes(nil) != 1 {
			t.Fail()
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		if sizeOfBytes([]byte{}) != 1 {
			t.Fail()
		}
	})

	t.Run("small <24 byte slice", func(t *testing.T) {
		if sizeOfBytes(make([]byte, 23)) != 1+23 {
			t.Fail()
		}
	})
	t.Run("small <256 byte slice", func(t *testing.T) {
		if sizeOfBytes(make([]byte, 200)) != 1+1+200 {
			t.Fail()
		}

	})
	t.Run("<64K byte slice", func(t *testing.T) {
		if sizeOfBytes(make([]byte, 1200)) != 1+2+1200 {
			t.Fail()
		}
	})

	t.Run(">64K byte slice", func(t *testing.T) {
		if sizeOfBytes(make([]byte, 70000)) != 1+4+70000 {
			t.Fail()
		}
	})
}

func TestEncodeBytesWithType(t *testing.T) {
}
