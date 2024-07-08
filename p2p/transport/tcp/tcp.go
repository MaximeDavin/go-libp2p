package tcp

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/transport"
)

type Option func(*TcpTransport) error

// TcpTransport is the TCP transport.
type TcpTransport struct{}

func NewTCPTransport(upgrader transport.Upgrader, rcmgr network.ResourceManager, opts ...Option) (*TcpTransport, error) {
}
