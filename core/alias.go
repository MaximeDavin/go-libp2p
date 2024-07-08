package core

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// ProtocolID aliases protocol.ID.
//
// Refer to the docs on that type for more info.
type ProtocolID = protocol.ID

// Stream aliases network.Stream.
//
// Refer to the docs on that type for more info.
type Stream = network.Stream

// PeerID aliases peer.ID.
//
// Refer to the docs on that type for more info.
type PeerID = peer.ID
