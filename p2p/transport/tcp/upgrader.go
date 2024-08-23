package tcp

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/transport"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	manet "github.com/multiformats/go-multiaddr/net"
	mss "github.com/multiformats/go-multistream"
)

const YAMUX_ID = "/yamux/1.0.0"

type conn struct {
	network.MuxedConn
	network.SecureConnMixin
	network.ConnMultiaddrsMixin
}

var _ transport.UpgradedConn = &conn{}

func upgrade(ctx context.Context, t *TcpTransport, rawConn manet.Conn, pid peer.ID, direction network.Direction) (transport.UpgradedConn, error) {
	// Upgrade security
	_, err := mss.SelectOneOf(t.SecuritySupported, rawConn)
	if err != nil {
		// TODO: error "Not compatible with noise"
		return nil, err
	}
	sconn, err := noise.Secure(ctx, rawConn, direction, pid)
	if err != nil {
		return nil, err
	}
	log.Printf("1")
	// Upgrade stream muxer
	streamproto, err := mss.SelectOneOf(t.MuxSupported, sconn)
	if err != nil {
		// Muxer negociation failed
		return nil, err
	}
	log.Printf("2")

	if streamproto != YAMUX_ID {
		// TODO: error handling
		return nil, errors.New(": this muxer is not supported")
	}
	log.Printf("3")
	mconn, err := yamux.Multiplex(sconn, direction)
	tc := &conn{
		MuxedConn:           mconn,
		SecureConnMixin:     sconn,
		ConnMultiaddrsMixin: rawConn,
	}
	return tc, nil
}
