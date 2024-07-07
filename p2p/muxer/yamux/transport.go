package yamux

import "github.com/libp2p/go-yamux/v4"

var DefaultTransport *Transport

const ID = "/yamux/1.0.0"

// Transport implements mux.Multiplexer that constructs
// yamux-backed muxed connections.
type Transport yamux.Config
