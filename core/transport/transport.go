package transport

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// The Transport interface allows you to open connections to other peers
// by dialing them, and also lets you listen for incoming connections.
//
// Connections returned by Dial and passed into Listeners are of type
// CapableConn, which means that they have been upgraded to support
// stream multiplexing and connection security (encryption and authentication).
type Transport interface {
	Dial(ctx context.Context, addr ma.Multiaddr, pid peer.ID) (UpgradedConn, error)
	Listen(addr ma.Multiaddr) (Listener, error)
}

type Listener interface {
	Accept() (UpgradedConn, error)
	Close() error
	Multiaddr() ma.Multiaddr
}

// A UpgradedConn represents a connection that has been secured and multiplexed
type UpgradedConn interface {
	network.MuxedConn
	network.SecureConnMixin
	network.ConnMultiaddrsMixin
}
