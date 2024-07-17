package transport

import (
	"context"

	"github.com/hashicorp/yamux"
	mplex "github.com/libp2p/go-mplex"
)

type YamuxAdapter struct {
	Session *yamux.Session
}

func (y *YamuxAdapter) OpenStream(ctx context.Context) (Stream, error) {
	stream, err := y.Session.OpenStream()
	return Stream(stream), err
}

func (y *YamuxAdapter) AcceptStream() (Stream, error) {
	stream, err := y.Session.AcceptStream()
	return Stream(stream), err
}

func (y *YamuxAdapter) Close() error {
	return y.Session.Close()
}

type MplexAdapter struct {
	Session *mplex.Multiplex
}

func (m *MplexAdapter) OpenStream(ctx context.Context) (Stream, error) {
	stream, err := m.Session.NewStream(ctx)
	return Stream(stream), err
}

func (m *MplexAdapter) AcceptStream() (Stream, error) {
	stream, err := m.Session.Accept()
	return Stream(stream), err
}

func (m *MplexAdapter) Close() error {
	return (m.Session.Close())
}
