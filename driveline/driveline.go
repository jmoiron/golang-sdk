/*
Package driveline is a client library for Driveline.


import "github.com/1533-systems/golang-sdk/driveline"

func main() {
}

*/
package driveline

import (
	"errors"

	ws "github.com/1533-systems/golang-sdk/driveline/websocket"
)

// ErrMaxReconnect indicates that the client has reached its maximum number of connection attempts.
var ErrMaxReconnect = ws.ErrMaxReconnect

// ErrClosed indicates that the connection is closed.
var ErrClosed = errors.New("connection closed")

// ErrInvalidServerMessage indicates that the client cannot decode server messages.
var ErrInvalidServerMessage = errors.New("invalid server message")
