package blank

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

// BlankHost is the thinnest implementation of the host.Host interface
type BlankHost struct {
}

type config struct{}

type Option = func(cfg *config)

func NewBlankHost(n network.Network, options ...Option) *BlankHost {}

func (bh *BlankHost) NewStream(ctx context.Context, p peer.ID, protos ...protocol.ID) (network.Stream, error) {
}

func (bh *BlankHost) ID() peer.ID {}

func (bh *BlankHost) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {}

func (bh *BlankHost) Connect(ctx context.Context, ai peer.AddrInfo) error {}

// TODO: also not sure this fits... Might be better ways around this (leaky abstractions)
func (bh *BlankHost) Network() network.Network {}

func (bh *BlankHost) Peerstore() peerstore.Peerstore {}

func (bh *BlankHost) Addrs() []ma.Multiaddr {}

func (bh *BlankHost) Close() error {}

// TODO: i'm not sure this really needs to be here
func (bh *BlankHost) Mux() protocol.Switch {}

func (bh *BlankHost) RemoveStreamHandler(pid protocol.ID) {}
