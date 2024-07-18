package noise

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	manet "github.com/multiformats/go-multiaddr/net"
)

// ID is the protocol ID for noise
const ID = "/noise"

type Transport struct {
	protocolID protocol.ID
	localID    peer.ID
	privateKey crypto.PrivKey
	muxers     []protocol.ID
}

func Secure(conn manet.Conn, initiator bool, pid string) (manet.Conn, error) {
	return conn, nil
}
