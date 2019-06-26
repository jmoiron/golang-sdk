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
	"time"
)

type optStoreOption uint16

const (
	optStoreTTLOption = optStoreOption(1 << iota)
	optStoreCASOption
)

// StoreOptions is used to configure  the behavior of the Store operation.
// It can be used to set the TTL of the record and also apply
// CAS (compare and swap) when the record is stored.
type StoreOptions struct {
	assigned    optStoreOption
	ttl         uint64 // in milliseconds
	casRecordID RecordID
}

// WithTTL sets the TTL of the record
func (o *StoreOptions) WithTTL(d time.Duration) *StoreOptions {
	if o != nil {
		o.assigned |= optStoreTTLOption
		o.ttl = uint64(d / time.Millisecond)
	}
	return o
}

// CompareAndSwap sets the expected RecordID used when storing a record.d
func (o *StoreOptions) CompareAndSwap(recordID RecordID) *StoreOptions {
	if o != nil {
		o.assigned |= optStoreCASOption
		o.casRecordID = recordID
	}
	return o
}
