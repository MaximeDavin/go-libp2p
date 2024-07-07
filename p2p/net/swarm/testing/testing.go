package testing

import (
	"testing"

	"github.com/libp2p/go-libp2p/p2p/net/swarm"
)

type config struct {
	disableQUIC bool
}

// Option is an option that can be passed when constructing a test swarm.
type Option func(*testing.T, *config)

// OptDisableQUIC disables QUIC.
var OptDisableQUIC Option = func(_ *testing.T, c *config) {
	c.disableQUIC = true
}

// GenSwarm generates a new test swarm.
func GenSwarm(t *testing.T, opts ...Option) *swarm.Swarm {}
