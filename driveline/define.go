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
	"sync"
)

const maxAliases = 256

type defines struct {
	aliases     []uint8
	aliasMap    map[uint8]string
	freeAliases int
	mu          sync.Mutex
}

func (d *defines) reset() {
	d.mu.Lock()
	d.aliases = make([]byte, maxAliases)
	for i := 0; i < maxAliases; i++ {
		d.aliases[i] = uint8(maxAliases - i)
	}
	d.aliasMap = make(map[uint8]string, maxAliases)
	d.freeAliases = maxAliases
	d.mu.Unlock()
}

func (d *defines) allocate(name string) streamID {
	d.mu.Lock()
	if d.freeAliases <= 0 {
		d.mu.Unlock()
		return textualStreamID(name)
	}
	d.freeAliases -= 1
	id := d.aliases[d.freeAliases]
	d.aliasMap[id] = name
	d.mu.Unlock()
	return numericStreamID(id)
}

func (d *defines) release(id streamID) {
	if id.isNumeric() {
		d.mu.Lock()
		numID := byte(id.numericID())
		d.aliases[d.freeAliases] = numID
		d.freeAliases += 1
		delete(d.aliasMap, numID)
		d.mu.Unlock()
	}
}
