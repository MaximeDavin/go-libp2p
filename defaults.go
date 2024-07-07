package libp2p

import "github.com/libp2p/go-libp2p/p2p/muxer/yamux"

// DefaultMuxers configures libp2p to use the stream connection multiplexers.
//
// libp2p instead of replacing them.
var DefaultMuxers = Muxer(yamux.ID, yamux.DefaultTransport)

// FallbackDefaults applies default options to the libp2p node if and only if no
// other relevant options have been applied. will be appended to the options
// passed into New.
var FallbackDefaults Option = func(cfg *Config) error {}
