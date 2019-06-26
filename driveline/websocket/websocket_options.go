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

package websocket

import (
	"net/http"
	"time"
)

type Option func(options *webSocketOptions)

type webSocketOptions struct {
	maxReconnect      int
	reconnectWait     time.Duration
	maxReconnectWait  time.Duration
	maxInFlight       int
	connectHandler    ConnectHandler
	disconnectHandler DisconnectHandler
	messageHandler    MessageHandler
	failureHandler    FailureHandler
	errorHandler      func(error)
	connectTimeout    time.Duration
	httpHeaders       http.Header
	httpClient        *http.Client
}

func (o *webSocketOptions) configure(options []Option) {
	o.connectTimeout = defaultConnectTimeout
	o.maxReconnect = -1
	o.reconnectWait = defaultReconnectWait
	o.maxReconnectWait = defaultMaxReconnectWait
	o.maxInFlight = defaultMaxInFlight
	o.connectHandler = func(WebSocket) {}
	o.disconnectHandler = func(WebSocket) {}
	o.messageHandler = func(WebSocket, []byte) {}
	o.failureHandler = func(WebSocket, error) {}
	o.errorHandler = func(error) {}
	o.configureHTTP()
	for _, configure := range options {
		configure(o)
	}
}

func MaxReconnect(max int) Option {
	return func(ws *webSocketOptions) {
		ws.maxReconnect = max
	}
}

func ReconnectForever() Option {
	return func(ws *webSocketOptions) {
		ws.maxReconnect = -1
	}
}

func ReconnectWait(d time.Duration) Option {
	return func(ws *webSocketOptions) {
		ws.reconnectWait = d
		if ws.maxReconnectWait < d {
			ws.maxReconnectWait = d * maxReconnectWaitRatio
		}
	}
}

func MaxReconnectWait(d time.Duration) Option {
	return func(ws *webSocketOptions) {
		ws.maxReconnectWait = d
	}

}

func OnMessage(handler func(WebSocket, []byte)) Option {
	return func(ws *webSocketOptions) {
		ws.messageHandler = handler
	}
}

func OnConnect(handler func(WebSocket)) Option {
	return func(ws *webSocketOptions) {
		ws.connectHandler = handler
	}
}

func OnDisconnect(handler func(WebSocket)) Option {
	return func(ws *webSocketOptions) {
		ws.disconnectHandler = handler
	}
}

func OnFailure(handler func(WebSocket, error)) Option {
	return func(ws *webSocketOptions) {
		ws.failureHandler = handler
	}
}

func MaxInFlight(count int) Option {
	return func(ws *webSocketOptions) {
		ws.maxInFlight = count
	}
}

func ErrorHandler(handler func(error)) Option {
	return func(ws *webSocketOptions) {
		ws.errorHandler = handler
	}
}
