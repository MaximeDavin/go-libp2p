package tcp

import (
	"context"
	"errors"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/transport"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	ma "github.com/multiformats/go-multiaddr"
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

func upgrade(ctx context.Context, rawConn manet.Conn, addr ma.Multiaddr, pid peer.ID, direction network.Direction) (transport.UpgradedConn, error) {
	// Upgrade security
	_, err := mss.SelectOneOf([]string{noise.ID}, rawConn)
	if err != nil {
		// TODO: error "Not compatible with noise"
		return nil, err
	}
	sconn, err := noise.Secure(ctx, rawConn, direction, pid)
	if err != nil {
		return nil, err
	}

	// Upgrade stream muxer
	var streamMuxerIDs = []string{YAMUX_ID}
	streamproto, err := mss.SelectOneOf(streamMuxerIDs, sconn)
	if err != nil {
		// Muxer negociation failed
		return nil, err
	}
	if streamproto != YAMUX_ID {
		// TODO: error handling
		return nil, errors.New("Programming error: this muxer is not supported")
	}
	mconn, err := yamux.Multiplex(sconn, direction)
	tc := &conn{
		MuxedConn:           mconn,
		SecureConnMixin:     sconn,
		ConnMultiaddrsMixin: rawConn,
	}
	return tc, nil
}
