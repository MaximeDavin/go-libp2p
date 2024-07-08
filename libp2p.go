package libp2p

import (
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
)

// Config describes a set of settings for a libp2p node.
type Config = config.Config

type Option = config.Option

func New(opts ...Option) (host.Host, error) {
}
