package network

import (
	"io"
	"time"

	"github.com/libp2p/go-libp2p/core/protocol"
)

// Stream represents a bidirectional channel between two agents in
// a libp2p network. "agent" is as granular as desired, potentially
// being a "request -> reply" pair, or whole protocols.
//
// Streams are backed by a multiplexer underneath the hood.
type Stream interface {
	io.Reader
	io.Writer
	io.Closer

	Protocol() protocol.ID

	// CloseWrite closes the stream for writing but leaves it open for
	// reading.
	//
	// CloseWrite does not free the stream, users must still call Close or
	// Reset.
	CloseWrite() error

	// Stat returns metadata pertaining to this stream.
	Stat() Stats

	// Conn returns the connection this stream is part of.
	Conn() Conn

	// Reset closes both ends of the stream. Use this to tell the remote
	// side to hang up and go away.
	Reset() error

	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}
