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
	"context"
	"time"

	ws "github.com/1533-systems/golang-sdk/driveline/websocket"
)

type clientOptions struct {
	client       *Client
	wsOptions    []ws.Option
	newWebSocket func(context.Context, string, ...ws.Option) (ws.WebSocket, error)
}

type option func(*clientOptions)

// ErrorHandler sets a global error handler. This handler can be used for logging.
func ErrorHandler(handler func(error)) option {
	return func(opts *clientOptions) {
		opts.client.errorHandler = handler
	}
}

// MaxReconnect limits the number of consecutive connection failure.
func MaxReconnect(max int) option {
	return func(opts *clientOptions) {
		opts.wsOptions = append(opts.wsOptions, ws.MaxReconnect(max))
	}
}

// ReconnectForever forces the client to try to reconnect forever.
func ReconnectForever() option {
	return func(opts *clientOptions) {
		opts.wsOptions = append(opts.wsOptions, ws.ReconnectForever())
	}
}

// ReconnectWait sets the minimum wait-period before the client tries to reconnect.
func ReconnectWait(d time.Duration) option {
	return func(opts *clientOptions) {
		opts.wsOptions = append(opts.wsOptions, ws.ReconnectWait(d))
	}
}

// MaxInFlight sets the maximum number of messages that can be buffered before they are sent to Driveline.
func MaxInFlight(count int) option {
	return func(opts *clientOptions) {
		opts.wsOptions = append(opts.wsOptions, ws.MaxInFlight(count))
	}
}

// for testing
func websocketProvider(provider func(context.Context, string, ...ws.Option) (ws.WebSocket, error)) option {
	return func(opts *clientOptions) {
		opts.newWebSocket = provider
	}
}
