package network

import (
	"errors"
	"io"
	"time"
)

type MuxedStream interface {
	io.Reader
	io.Writer
	io.Closer
	// CloseWrite() error
	// Reset() error

	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

// Multiplexer wraps a net.Conn with a stream multiplexing
// implementation and returns a MuxedConn that supports opening
// multiple streams over the underlying net.Conn
type Multiplexer interface{}

// ErrReset is returned when reading or writing on a reset stream.
var ErrReset = errors.New("stream reset")
