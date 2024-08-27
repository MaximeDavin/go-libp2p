package yamux

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"

	"github.com/hashicorp/yamux"
)

var DefaultTransport *Transport

const ID = "/yamux/1.0.0"

// Transport implements mux.Multiplexer that constructs
// yamux-backed muxed connections.
type Transport yamux.Config

func Multiplex(conn network.SecureConn, direction network.Direction) (network.MuxedConn, error) {
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, err
	}
	return &yamuxConn{Session: session}, nil
}

type yamuxConn struct {
	*yamux.Session
}

var _ network.MuxedConn = &yamuxConn{}

func (y *yamuxConn) OpenStream(ctx context.Context) (network.MuxedStream, error) {
	stream, err := y.Session.OpenStream()
	return network.MuxedStream(stream), err
}

func (y *yamuxConn) AcceptStream() (network.MuxedStream, error) {
	stream, err := y.Session.AcceptStream()
	return network.MuxedStream(stream), err
}

func (y *yamuxConn) Close() error {
	return y.Session.Close()
}

func (y *yamuxConn) IsClosed() bool {
	return y.Session.IsClosed()
}
