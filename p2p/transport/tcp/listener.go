package tcp

import (
	"github.com/libp2p/go-libp2p/core/transport"
	manet "github.com/multiformats/go-multiaddr/net"
)

type TcpListener struct {
	manet.Listener
	// Every accepted and upgraded connection in the listen loop's
	// goroutines are sent to this channel
	incoming chan transport.UpgradedConn
	// Store error happening during listen loop
	errs chan error
}

var _ transport.Listener = &TcpListener{}

func NewTCPListener() (*TcpTransport, error) {
	return &TcpTransport{}, nil
}

// Accepts incoming connections from the listen loop
func (l *TcpListener) Accept() (transport.UpgradedConn, error) {
	select {
	case c := <-l.incoming:
		return c, nil
	case err := <-l.errs:
		log.Errorf("%s", err)
		return nil, err
	}
}

func (l *TcpListener) Close() error {
	for c := range l.incoming {
		c.Close()
	}
	return l.Listener.Close()
}
