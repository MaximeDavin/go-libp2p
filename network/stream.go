package network

import (
	"context"
	"io"

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
	SetProtocol(id protocol.ID) error
}

// Connection that support stream multiplexing
type StreamConn interface {
	io.Closer
	// OpenStream creates a new stream.
	OpenStream(context.Context) (Stream, error)

	// AcceptStream accepts a stream opened by the other side.
	AcceptStream() (Stream, error)
}
