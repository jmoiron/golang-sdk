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

type streamID interface {
	isNumeric() bool
	numericID() uint64
	textualID() string
}

type numericStreamID int

var _ streamID = numericStreamID(0)
var _ streamID = textualStreamID("")

func (numericStreamID) isNumeric() bool {
	return true
}

func (id numericStreamID) numericID() uint64 {
	return uint64(id)
}

func (id numericStreamID) textualID() string {
	return ""
}

type textualStreamID string

func (textualStreamID) isNumeric() bool {
	return false
}

func (id textualStreamID) numericID() uint64 {
	return 0
}

func (id textualStreamID) textualID() string {
	return string(id)
}
