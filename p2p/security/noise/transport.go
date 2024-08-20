package noise

import (
	"context"
	"net"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

// ID is the protocol ID for noise
const ID = "/noise"

func New() {}

type secureConn struct {
	insecureConn manet.Conn
	localID      peer.ID
	remoteID     peer.ID
}

func Secure(ctx context.Context, insecureConn manet.Conn, direction network.Direction, pid peer.ID) (network.SecureConn, error) {
	sconn := &secureConn{
		insecureConn: insecureConn,
		localID:      "",
		remoteID:     pid,
	}
	return sconn, nil
}

func (s *secureConn) LocalAddr() net.Addr {
	return s.insecureConn.LocalAddr()
}

func (s *secureConn) LocalMultiaddr() multiaddr.Multiaddr {
	return s.insecureConn.LocalMultiaddr()
}

func (s *secureConn) RemoteAddr() net.Addr {
	return s.insecureConn.RemoteAddr()
}

func (s *secureConn) RemoteMultiaddr() multiaddr.Multiaddr {
	return s.insecureConn.RemoteMultiaddr()
}

func (s *secureConn) SetDeadline(t time.Time) error {
	return s.insecureConn.SetDeadline(t)
}

func (s *secureConn) SetReadDeadline(t time.Time) error {
	return s.insecureConn.SetReadDeadline(t)
}

func (s *secureConn) SetWriteDeadline(t time.Time) error {
	return s.insecureConn.SetWriteDeadline(t)
}

func (s *secureConn) Close() error {
	return s.insecureConn.Close()
}

func (s *secureConn) LocalPeer() peer.ID {
	return s.localID
}

func (s *secureConn) RemotePeer() peer.ID {
	return s.remoteID
}

func (s *secureConn) Read(buf []byte) (int, error) {
	// TODO: implements secured read
	return s.insecureConn.Read(buf)
}

func (s *secureConn) Write(data []byte) (int, error) {
	// TODO: implements secured write
	return s.insecureConn.Write(data)
}
