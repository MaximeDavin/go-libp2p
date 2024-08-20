package tcp

import (
	"errors"

	"github.com/libp2p/go-libp2p/core/transport"
	manet "github.com/multiformats/go-multiaddr/net"
)

type TcpListener struct {
	manet.Listener
	incoming chan transport.UpgradedConn
}

var _ transport.Listener = &TcpListener{}

func NewTCPListener() (*TcpTransport, error) {
	return &TcpTransport{}, nil
}

func (l *TcpListener) Accept() (transport.UpgradedConn, error) {
	for c := range l.incoming {
		return c, nil

	}
	return nil, errors.New("inconsistent length for security transports")
}

func (l *TcpListener) Close() error {
	return l.Listener.Close()
}
