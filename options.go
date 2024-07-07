package libp2p

import (
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

// ListenAddrs configures libp2p to listen on the given addresses.
func ListenAddrs(addrs ...ma.Multiaddr) Option {}

// UserAgent sets the libp2p user-agent sent along with the identify protocol
func UserAgent(userAgent string) Option {}

// ConnectionGater configures libp2p to use the given ConnectionGater
// to actively reject inbound/outbound connections based on the lifecycle stage
// of the connection.
//
// For more information, refer to go-libp2p/core.ConnectionGater.
func ConnectionGater(cg connmgr.ConnectionGater) Option {}

// Transport configures libp2p to use the given transport (or transport
// constructor).
//
// The transport can be a constructed transport.Transport or a function taking
// any subset of this libp2p node's:
// * Transport Upgrader (*tptu.Upgrader)
// * Host
// * Stream muxer (muxer.Transport)
// * Security transport (security.Transport)
// * Private network protector (pnet.Protector)
// * Peer ID
// * Private Key
// * Public Key
// * Address filter (filter.Filter)
// * Peerstore
func Transport(constructor interface{}, opts ...interface{}) Option {}

// Muxer configures libp2p to use the given stream multiplexer.
// name is the protocol name.
func Muxer(name string, muxer network.Multiplexer) Option {}

// Security configures libp2p to use the given security transport (or transport
// constructor).
//
// Name is the protocol name.
//
// The transport can be a constructed security.Transport or a function taking
// any subset of this libp2p node's:
// * Public key
// * Private key
// * Peer ID
// * Host
// * Network
// * Peerstore
func Security(name string, constructor interface{}) Option {}

// Ping will configure libp2p to support the ping service; enable by default.
func Ping(enable bool) Option {}

// NATPortMap configures libp2p to use the default NATManager. The default
// NATManager will attempt to open a port in your network's firewall using UPnP.
func NATPortMap() Option {}

// AddrsFactory configures libp2p to use the given address factory.
func AddrsFactory(factory config.AddrsFactory) Option {}

// DisableRelay configures libp2p to disable the relay transport.
func DisableRelay() Option {}

// ResourceManager configures libp2p to use the given ResourceManager.
// When using the p2p/host/resource-manager implementation of the ResourceManager interface,
// it is recommended to set limits for libp2p protocol by calling SetDefaultServiceLimits.
func ResourceManager(rcmgr network.ResourceManager) Option {}

// Identity configures libp2p to use the given private key to identify itself.
func Identity(sk crypto.PrivKey) Option {}
