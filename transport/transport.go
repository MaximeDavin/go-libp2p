package transport

import (
	"context"
	"io"

	ma "github.com/multiformats/go-multiaddr"
)

// Transport represents any device by which you can connect to and accept
// connections from other peers.
//
// The Transport interface allows you to open connections to other peers
// by dialing them, and also lets you listen for incoming connections.
// This connections are secured and multiplexed

type Transport interface {
	// Dial dials a remote peer.
	Dial(ctx context.Context, addr ma.Multiaddr, pid string) (StreamConn, error)
	// Listen listens on the passed multiaddr.
	Listen(laddr ma.Multiaddr) (Listener, error)
}

type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (StreamConn, error)
	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error
}

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}

// Connection that support stream multiplexing
type StreamConn interface {
	io.Closer
	// OpenStream creates a new stream.
	OpenStream(context.Context) (Stream, error)

	// AcceptStream accepts a stream opened by the other side.
	AcceptStream() (Stream, error)
}
