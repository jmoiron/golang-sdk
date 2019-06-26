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
	"errors"
)

type serverMsg struct {
	consumerID uint64
	err        error
	records    []Record
}

func decodeServerMessage(buf []byte) (*serverMsg, error) {
	if !isArray(buf[0]) {
		return nil, ErrInvalidServerMessage
	}
	itemCount, buf, err := decodeNumber(buf)
	if err != nil {
		return nil, err
	}
	_ = buf[1] // bounds check elimination
	if buf[0] == (cborTextString|4) && buf[1] == 'd' {
		return decodeDataMessage(buf[5:], int(itemCount)-3)
	} else if buf[0] == (cborTextString | 3) {
		if buf[1] == 's' {
			return decodeSyncMessage(buf[4:])
		} else if buf[1] == 'e' {
			return decodeErrorMessage(buf[4:])
		}
	}
	return nil, ErrInvalidServerMessage
}

func decodeDataMessage(buf []byte, recordCount int) (*serverMsg, error) {
	if !isUnsignedInteger(buf[0]) {
		return nil, ErrInvalidServerMessage
	}
	consumerID, buf, err := decodeNumber(buf)
	if err != nil {
		return nil, err
	}
	records := make([]Record, recordCount)
	d := buf[0]
	switch {
	case isArray(d):
		var tagCount uint64
		tagCount, buf, err = decodeNumber(buf)
		if err != nil {
			return nil, err
		}
		// Must be an even number of entries
		if tagCount%2 != 0 {
			return nil, ErrInvalidServerMessage
		}
		for i := 0; i < int(tagCount); i += 2 {
			d := buf[0]
			if !isUnsignedInteger(d) {
				return nil, ErrInvalidServerMessage
			}
			if d != encodedMessageIdTag {
				return nil, ErrInvalidServerMessage
			}
			if !isArray(buf[1]) {
				return nil, ErrInvalidServerMessage
			}
			var idCnt uint64
			idCnt, buf, err = decodeNumber(buf[1:])
			if err != nil {
				return nil, err
			}
			if int(idCnt) != recordCount {
				return nil, ErrInvalidServerMessage
			}
			for i := 0; i < recordCount; i++ {
				records[i].RecordID, buf, err = decodeRecordID(buf)
				if err != nil {
					return nil, err
				}
			}
		}
	case isBlank(d):
		buf = buf[1:]
	default:
		return nil, ErrInvalidServerMessage
	}

	if recordCount == 1 && isUndefined(buf[0]) {
		return &serverMsg{consumerID: consumerID}, nil
	}
	for i := 0; i < recordCount; i++ {
		records[i].Record, buf, err = decodeBytes(buf)
	}
	return &serverMsg{consumerID: consumerID, records: records}, nil
}

func decodeErrorMessage(data []byte) (*serverMsg, error) {
	if !isUnsignedInteger(data[0]) {
		return nil, ErrInvalidServerMessage
	}
	consumerID, data, err := decodeNumber(data)
	if err != nil {
		return nil, err
	}
	errorMsg, data, err := decodeString(data)
	if err != nil {
		return nil, err
	}
	return &serverMsg{consumerID: consumerID, err: errors.New(errorMsg)}, nil
}

func decodeSyncMessage(buf []byte) (*serverMsg, error) {
	if !isUnsignedInteger(buf[0]) {
		return nil, ErrInvalidServerMessage
	}
	consumerID, buf, err := decodeNumber(buf)
	if err != nil {
		return nil, err
	}
	return &serverMsg{consumerID: consumerID}, nil
}

func decodeNumber(buf []byte) (uint64, []byte, error) {
	size := lenCode(buf[0])
	switch {
	case size < 24:
		return size, buf[1:], nil
	case size == 24:
		return uint64(buf[1]), buf[2:], nil
	case size == 25:
		result := uint64(binary.BigEndian.Uint16(buf[1:]))
		return result, buf[3:], nil
	case size == 26:
		result := uint64(binary.BigEndian.Uint32(buf[1:]))
		return result, buf[5:], nil
	case size == 27:
		result := uint64(binary.BigEndian.Uint64(buf[1:]))
		return result, buf[9:], nil
	default:
		return 0, nil, ErrInvalidServerMessage
	}
}

func decodeBytes(buf []byte) ([]byte, []byte, error) {
	d := buf[0]
	if isBlank(d) {
		return nil, buf[1:], nil
	}
	if !isByteString(d) {
		return nil, nil, ErrInvalidServerMessage
	}
	buf = buf[1:]
	size := lenCode(d)
	switch {
	case size < 24:
		// pass
	case size == 24:
		size = uint64(buf[0])
		buf = buf[1:]
	case size == 25:
		size = uint64(binary.BigEndian.Uint16(buf))
		buf = buf[2:]
	case size == 26:
		size = uint64(binary.BigEndian.Uint32(buf))
		buf = buf[4:]
	case size == 27:
		size = binary.BigEndian.Uint64(buf)
		buf = buf[8:]
	default:
		return nil, nil, ErrInvalidServerMessage
	}
	return buf[:size], buf[size:], nil
}

func decodeString(buf []byte) (string, []byte, error) {
	d := buf[0]
	if isBlank(d) {
		return "", buf, nil
	}
	if !isTextString(d) {
		return "", nil, ErrInvalidServerMessage
	}
	buf = buf[1:]
	size := lenCode(d)
	switch {
	case size < 24:
		// pass
	case size == 24:
		size = uint64(buf[0])
		buf = buf[1:]
	case size == 25:
		size = uint64(binary.BigEndian.Uint16(buf))
		buf = buf[2:]
	case size == 26:
		size = uint64(binary.BigEndian.Uint32(buf))
		buf = buf[4:]
	case size == 27:
		size = binary.BigEndian.Uint64(buf)
		buf = buf[8:]
	default:
		return "", nil, ErrInvalidServerMessage
	}
	return string(buf[:size]), buf[size:], nil
}

func decodeRecordID(buf []byte) (RecordID, []byte, error) {
	id, buf, err := decodeBytes(buf)
	return RecordID(id), buf, err
}
