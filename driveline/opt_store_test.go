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
	"reflect"
	"testing"
	"time"
)

func TestStoreOptions_CompareAndSwap(t *testing.T) {
	var o StoreOptions
	o.CompareAndSwap(testRecordID)
	if !reflect.DeepEqual(testRecordID, o.casRecordID) {
		t.Fail()
	}
	if o.assigned&optStoreCASOption == 0 {
		t.Fail()
	}
}

func TestStoreOptions_WithTTL(t *testing.T) {
	var o StoreOptions
	o.WithTTL(2 * time.Second)

	if o.ttl != 2000 {
		t.Fail()
	}
	if o.assigned&optStoreTTLOption == 0 {
		t.Fail()
	}

}
