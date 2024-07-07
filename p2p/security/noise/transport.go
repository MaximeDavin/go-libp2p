package noise

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	tptu "github.com/libp2p/go-libp2p/p2p/net/upgrader"
)

// ID is the protocol ID for noise
const ID = "/noise"

type Transport struct {
	protocolID protocol.ID
	localID    peer.ID
	privateKey crypto.PrivKey
	muxers     []protocol.ID
}

// New creates a new Noise transport using the given private key as its
// libp2p identity key.
func New(id protocol.ID, privkey crypto.PrivKey, muxers []tptu.StreamMuxer) (*Transport, error) {}
