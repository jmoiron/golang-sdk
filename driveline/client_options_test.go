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
	"errors"
	"testing"
	"time"
)

func TestClientOptions(t *testing.T) {

	t.Run("sets an error handler", func(t *testing.T) {
		opts := newClientOptions()

		var errs []error
		h := func(err error) {
			errs = append(errs, err)
		}
		ErrorHandler(h)(&opts)
		opts.client.errorHandler(errors.New("ouch"))

		if len(errs) != 1 {
			t.Fail()
		}
		if errs[0].Error() != "ouch" {
			t.Fail()
		}
	})

	t.Run("sets an error handler", func(t *testing.T) {
		opts := newClientOptions()

		var errs []error
		h := func(err error) {
			errs = append(errs, err)
		}
		ErrorHandler(h)(&opts)
		opts.client.errorHandler(errors.New("ouch"))

		if len(errs) != 1 {
			t.Fail()
		}
		if errs[0].Error() != "ouch" {
			t.Fail()
		}
	})

	t.Run("sets max number of reconnect", func(t *testing.T) {
		opts := newClientOptions()
		if len(opts.wsOptions) != 0 {
			t.Fail()
		}
		MaxReconnect(123)(&opts)
		if len(opts.wsOptions) != 1 {
			t.Fail()
		}
	})

	t.Run("configures the WeboSocket options to reconnect forever", func(t *testing.T) {
		opts := newClientOptions()
		if len(opts.wsOptions) != 0 {
			t.Fail()
		}
		ReconnectForever()(&opts)
		if len(opts.wsOptions) != 1 {
			t.Fail()
		}
	})

	t.Run("configures the WebSocket wait time", func(t *testing.T) {
		opts := newClientOptions()
		if len(opts.wsOptions) != 0 {
			t.Fail()
		}
		ReconnectWait(time.Second)(&opts)
		if len(opts.wsOptions) != 1 {
			t.Fail()
		}
	})

	t.Run("configures the WebSocket maxInFlight", func(t *testing.T) {
		opts := newClientOptions()
		if len(opts.wsOptions) != 0 {
			t.Fail()
		}
		MaxInFlight(123)(&opts)
		if len(opts.wsOptions) != 1 {
			t.Fail()
		}
	})

}

func newClientOptions(options ...option) clientOptions {
	var opts clientOptions
	opts.client = new(Client)
	for _, option := range options {
		option(&opts)
	}
	return opts
}
