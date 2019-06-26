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

type optQueryOption uint16

const (
	optRecordQueryOption = optQueryOption(1 << iota)
)

// QueryOptions configures ContinuousQueryOptions or QueryOptions operations
type QueryOptions struct {
	assigned     optQueryOption
	fromRecordID RecordID
}

// FromStreamHead indicates that the query operation should start as far back
// as possible.
func (o *QueryOptions) FromStreamHead() *QueryOptions {
	if o != nil {
		o.assigned |= optRecordQueryOption
		o.fromRecordID = recordIDHead
		return o
	}
	return &QueryOptions{
		assigned:     optRecordQueryOption,
		fromRecordID: recordIDHead,
	}
}

// fromStreamTail indicates that the query operation should apply only for new
// records.
func (o *QueryOptions) fromStreamTail() *QueryOptions {
	if o != nil {
		o.assigned |= optRecordQueryOption
		o.fromRecordID = recordIDTail
		return o
	}

	return &QueryOptions{
		assigned:     optRecordQueryOption,
		fromRecordID: recordIDTail,
	}
}

// FromRecordID specifies the RecordID that to use query operation should use
// as a starting point.
func (o *QueryOptions) FromRecordID(id RecordID) *QueryOptions {
	if o != nil {
		o.assigned |= optRecordQueryOption
		o.fromRecordID = id
		return o
	}
	return &QueryOptions{
		assigned:     optRecordQueryOption,
		fromRecordID: id,
	}
}
