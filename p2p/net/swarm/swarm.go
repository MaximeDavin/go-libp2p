package swarm

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Swarm is a connection muxer, allowing connections to other peers to
// be opened and closed, while still using the same Chan for all
// communication. The Chan sends/receives Messages, which note the
// destination or source Peer.
type Swarm struct{}

// ClosePeer closes all connections to the given peer.
func (s *Swarm) ClosePeer(p peer.ID) error {}

// Connectedness returns our "connectedness" state with the given peer.
//
// To check if we have an open connection, use `s.Connectedness(p) ==
// network.Connected`.
func (s *Swarm) Connectedness(p peer.ID) network.Connectedness {}

// Conns returns a slice of all connections.
func (s *Swarm) Conns() []network.Conn {}

// Peers returns a copy of the set of peers swarm is connected to.
func (s *Swarm) Peers() []peer.ID {}

func (s *Swarm) ConnsToPeer(p peer.ID) []network.Conn {}

// Notify signs up Notifiee to receive signals when events happen
func (s *Swarm) Notify(f network.Notifiee) {}
