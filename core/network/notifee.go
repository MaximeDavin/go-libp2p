package network

import (
	ma "github.com/multiformats/go-multiaddr"
)

// Notifiee is an interface for an object wishing to receive
// notifications from a Network.
type Notifiee interface {
	// 	Listen(Network, ma.Multiaddr)      // called when network starts listening on an addr
	// 	ListenClose(Network, ma.Multiaddr) // called when network stops listening on an addr
	// 	Connected(Network, Conn)           // called when a connection opened
	// 	Disconnected(Network, Conn)        // called when a connection closed
}

// NotifyBundle implements Notifiee by calling any of the functions set on it,
// and nop'ing if they are unset. This is the easy way to register for
// notifications.
type NotifyBundle struct {
	ListenF      func(Network, ma.Multiaddr)
	ListenCloseF func(Network, ma.Multiaddr)

	ConnectedF    func(Network, Conn)
	DisconnectedF func(Network, Conn)
}
