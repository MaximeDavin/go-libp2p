package upgrader

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type StreamMuxer struct {
	ID    protocol.ID
	Muxer network.Multiplexer
}
