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
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/1533-systems/golang-sdk/driveline/bininfo"
)

var (
	ErrHandshake             = errors.New("invalid handshake")
	ErrInvalidProtocolScheme = errors.New("URL scheme must be ws or wss")
)

func (o *webSocketOptions) configureHTTP() {
	o.httpClient = http.DefaultClient
	o.httpHeaders = make(http.Header)
	o.httpHeaders.Set("Connection", "Upgrade")
	o.httpHeaders.Set("Upgrade", "websocket")
	o.httpHeaders.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	o.httpHeaders.Set("Sec-WebSocket-Protocol", "driveline")
	o.httpHeaders.Set("Sec-WebSocket-Version", "13")
	o.httpHeaders.Set("User-Agent", "driveline/"+bininfo.VERSION+" go")
}

func (ws *webSocket) connect(endpoint string) error {
	u, err := url.Parse(ws.endpoint)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	case "http", "https":
	// ignore
	default:
		return ErrInvalidProtocolScheme
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header = ws.httpHeaders
	ctx, _ := context.WithTimeout(context.Background(), ws.connectTimeout)
	req = req.WithContext(ctx)
	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return err
	}

	if err := verifyHandshake(req, resp); err != nil {
		return err
	}
	ok := true
	ws.cnx, ok = resp.Body.(io.ReadWriteCloser)
	if !ok {
		return ErrHandshake
	}
	return nil
}

func verifyHandshake(_ *http.Request, response *http.Response) error {
	if response.StatusCode != http.StatusSwitchingProtocols {
		return ErrHandshake
	}
	if strings.ToLower(response.Header.Get("Connection")) != "upgrade" {
		return ErrHandshake
	}
	if strings.ToLower(response.Header.Get("Upgrade")) != "websocket" {
		return ErrHandshake
	}
	return nil
}

func HTTPClient(c *http.Client) Option {
	return func(o *webSocketOptions) {
		o.httpClient = c
	}
}

func HTTPHeaders(header http.Header) Option {
	return func(o *webSocketOptions) {
		for k, v := range header {
			for _, val := range v {
				o.httpHeaders.Add(k, val)
			}
		}
	}
}

func ConnectTimeout(d time.Duration) Option {
	return func(ws *webSocketOptions) {
		ws.connectTimeout = d
	}
}
