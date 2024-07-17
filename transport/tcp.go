package transport

import (
	"context"
	"errors"

	"github.com/hashicorp/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"

	mplex "github.com/libp2p/go-mplex"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	mss "github.com/multiformats/go-multistream"
)

type TcpTransport struct{}

var _ Transport = &TcpTransport{}

func NewTCPTransport() (*TcpTransport, error) {
	return &TcpTransport{}, nil
}

type MuxedConn struct {
	yamux.Session
	mplex.Multiplex
}

func (t *TcpTransport) Dial(ctx context.Context, addr ma.Multiaddr, pid string) (StreamConn, error) {
	// Get raw connection
	var d manet.Dialer
	conn, err := d.DialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	const initiator = true

	// Upgrade security
	_, err = mss.SelectOneOf([]string{NOISE_ID}, conn)
	if err != nil {
		// Not compatible with noise
		return nil, err
	}
	sconn, err := noise.Secure(conn, initiator, pid)
	if err != nil {
		return nil, err
	}

	// Upgrade stream muxer
	var streamMuxerIDs = []string{YAMUX_ID, MPLEX_ID} // TODO: build this from config
	streamproto, err := mss.SelectOneOf(streamMuxerIDs, sconn)
	if err != nil {
		// No stream muxer compatible
		return nil, err
	}
	switch streamproto {
	case YAMUX_ID:
		// TODO: build and pass config
		yconn, err := yamux.Client(sconn, nil)
		if err != nil {
			return nil, err
		}
		yadapter := YamuxAdapter{Session: yconn}
		return &yadapter, nil
	case MPLEX_ID:
		mconn, err := mplex.NewMultiplex(sconn, initiator, nil)
		if err != nil {
			return nil, err
		}
		madapter := &MplexAdapter{Session: mconn}
		return madapter, nil
	default:
		return nil, errors.New("Programming error: this muxer is not supported")
	}

}

func (t *TcpTransport) Listen(addr ma.Multiaddr) (Listener, error) {
	return manet.Listen(addr)
}

// func (l *manet.Listener) Accept() (manet.Conn, error) {

// 	c, err := manet.Listener.Accept()
// }
