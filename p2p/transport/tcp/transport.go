package tcp

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/transport"

	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	mss "github.com/multiformats/go-multistream"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{
	"prefix":    "p2p",
	"transport": "TCP",
})

type TcpTransport struct {
	SecurityMuxer *mss.MultistreamMuxer[string]
	StreamMuxer   *mss.MultistreamMuxer[string]
	TcpTransportOptions
}

type TcpTransportOptions struct {
	// Maximum duration between accepting the raw TCP connection and returning
	// a fully upgraded connection
	AcceptTimeout time.Duration
	// ID of security protocol proposed in the protocol selection
	SecuritySupported string
	// ID of stream multiplexing protocols proposed in the protocol selection
	StreamSupported string
	// Function that upgrade a raw TCP connection into a fully upgraded connection
	// We will always use the one provided by default, except for some test cases
	// Upgrade func(ctx context.Context, t *TcpTransport, rawConn manet.Conn, pid peer.ID, direction network.Direction) (transport.UpgradedConn, error)
}

var DefaultTcpTransportOptions TcpTransportOptions = TcpTransportOptions{
	AcceptTimeout:     15 * time.Second,
	SecuritySupported: noise.ID,
	StreamSupported:   yamux.ID,
}

var _ transport.Transport = &TcpTransport{}

func NewTCPTransport() (*TcpTransport, error) {
	return NewTCPTransportWithOptions(DefaultTcpTransportOptions)
}

func NewTCPTransportWithOptions(options TcpTransportOptions) (*TcpTransport, error) {
	t := &TcpTransport{
		SecurityMuxer:       mss.NewMultistreamMuxer[string](),
		StreamMuxer:         mss.NewMultistreamMuxer[string](),
		TcpTransportOptions: options,
	}
	t.SecurityMuxer.AddHandler(options.SecuritySupported, nil)
	t.StreamMuxer.AddHandler(options.StreamSupported, nil)
	return t, nil
}

// Dial the given multiaddr and upgrade the resulting outbound connection
func (t *TcpTransport) Dial(ctx context.Context, addr ma.Multiaddr, pid peer.ID) (transport.UpgradedConn, error) {
	var d manet.Dialer
	conn, err := d.DialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return Upgrade(ctx, t, conn, pid, network.DirOutbound)
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
	l, err := manet.Listen(addr)
	if err != nil {
		return nil, err
	}

	incoming := make(chan transport.UpgradedConn)
	errs := make(chan error)

	go func() {
		defer func() {
			l.Close()
			close(incoming)
			close(errs)
		}()

		for {
			conn, err := l.Accept()
			if err != nil {
				// Temporary errors could be retried here
				// see https://go.dev/src/net/http/server.go?s=51504:51550#L3335
				// However it has been depecrated see https://groups.google.com/g/golang-nuts/c/-JcZzOkyqYI/m/xwaZzjCgAwAJ
				errs <- err
				return
			}
			log.Debugf("listener: connection from %s accepted",
				conn.RemoteMultiaddr(),
			)

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), t.AcceptTimeout)
				defer cancel()

				// Upgrade the connection: add security and stream multiplexing
				sconn, err := Upgrade(ctx, t, conn, "", network.DirInbound)
				if err != nil {
					log.Warnf("listener: connection from %s upgrade failed %s",
						conn.RemoteMultiaddr(),
						err,
					)
					errs <- err
					return
				}
				log.Debugf("listener: connection from %s updated",
					conn.RemoteMultiaddr(),
				)
				select {
				case incoming <- sconn:
				case <-ctx.Done():
					conn.Close()
				}
			}()

		}
	}()
	tcpl := &TcpListener{
		Listener: l,
		incoming: incoming,
		errs:     errs,
	}
	return tcpl, err

}
