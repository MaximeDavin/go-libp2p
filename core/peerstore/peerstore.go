package peerstore

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	ma "github.com/multiformats/go-multiaddr"
)

// Peerstore provides a thread-safe store of Peer related
// information.
type Peerstore interface {
	AddrBook
	PeerMetadata
	ProtoBook
	Metrics
	// PeerInfo returns a peer.PeerInfo struct for given peer.ID.
	// This is a small slice of the information Peerstore has on
	// that peer, useful to other services.
	PeerInfo(peer.ID) peer.AddrInfo

	// Peers returns all the peer IDs stored across all inner stores.
	Peers() peer.IDSlice
}

// AddrBook holds the multiaddrs of peers.
type AddrBook interface {
	// Addrs returns all known (and valid) addresses for a given peer.
	Addrs(p peer.ID) []ma.Multiaddr
}

// PeerMetadata can handle values of any type. Serializing values is
// up to the implementation. Dynamic type introspection may not be
// supported, in which case explicitly enlisting types in the
// serializer may be required.
//
// Refer to the docs of the underlying implementation for more
// information.
type PeerMetadata interface {
	// Get / Put is a simple registry for other peer-related key/value pairs.
	// If we find something we use often, it should become its own set of
	// methods. This is a last resort.
	Get(p peer.ID, key string) (interface{}, error)
	// Put(p peer.ID, key string, val interface{}) error

	// RemovePeer removes all values stored for a peer.
	// RemovePeer(peer.ID)
}

// Metrics tracks metrics across a set of peers.
type Metrics interface {
	// LatencyEWMA returns an exponentially-weighted moving avg.
	// of all measurements of a peer's latency.
	LatencyEWMA(peer.ID) time.Duration

	// RecordLatency records a new latency measurement
	RecordLatency(peer.ID, time.Duration)
}

// ProtoBook tracks the protocols supported by peers.
type ProtoBook interface {
	GetProtocols(peer.ID) ([]protocol.ID, error)
}
