package tcp

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/transport"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/sirupsen/logrus"
)

// TODO: config management
const DEFAULT_ACCEPT_TIMEOUT = 2 * time.Second

var log = logrus.WithFields(logrus.Fields{
	"prefix":    "p2p",
	"transport": "TCP",
})

type TcpTransport struct {
	acceptTimeout time.Duration
}

var _ transport.Transport = &TcpTransport{}

func NewTCPTransport() (*TcpTransport, error) {
	return &TcpTransport{
		acceptTimeout: DEFAULT_ACCEPT_TIMEOUT,
	}, nil
}

func (t *TcpTransport) Dial(ctx context.Context, addr ma.Multiaddr, pid peer.ID) (transport.UpgradedConn, error) {
	var d manet.Dialer
	conn, err := d.DialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return upgrade(ctx, conn, addr, pid, network.DirOutbound)
}

// Listen returns a listener that listens on the given multiaddr for inbound
// connections.
//
// For each connection accepted by the net.listener, a goroutine is run to
// upgrade the connection (adding security and stream multiplexing).
// Then, the upgraded connection is sent to the 'incoming' channel
// and received by the 'Accept' method of the returned listener.
// The use of goroutines allows the listener to keep accepting new connections
// while a connection is upgraded.
func (t *TcpTransport) Listen(addr ma.Multiaddr) (transport.Listener, error) {
	incoming := make(chan transport.UpgradedConn)

	l, err := manet.Listen(addr)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			l.Close()
			close(incoming)
		}()

		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			log.Debugf("listener: connection from %s accepted",
				conn.RemoteMultiaddr(),
			)

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), t.acceptTimeout)
				defer cancel()

				// Upgrade the connection: add security and stream multiplexing
				sconn, err := upgrade(ctx, conn, conn.RemoteMultiaddr(), "", network.DirInbound)
				if err != nil {
					log.Warnf("listener: connection from %s update failed %s",
						conn.RemoteMultiaddr(),
						err,
					)
					return
				}
				log.Debugf("listener: connection from %s updated",
					conn.RemoteMultiaddr(),
				)
				select {
				case incoming <- sconn:
				case <-ctx.Done():
					log.Warnf("")
					conn.Close()
				}
			}()

		}
	}()
	tcpl := &TcpListener{
		Listener: l,
		incoming: incoming,
	}
	return tcpl, err

}
