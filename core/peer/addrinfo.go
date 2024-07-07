package peer

import (
	ma "github.com/multiformats/go-multiaddr"
)

// AddrInfo is a small struct used to pass around a peer with
// a set of addresses (and later, keys?).
type AddrInfo struct {
	ID    ID
	Addrs []ma.Multiaddr
}

func (pi AddrInfo) String() string {}

// AddrInfosFromP2pAddrs converts a set of Multiaddrs to a set of AddrInfos.
func AddrInfosFromP2pAddrs(maddrs ...ma.Multiaddr) ([]AddrInfo, error) {}

// AddrInfoFromP2pAddr converts a Multiaddr to an AddrInfo.
func AddrInfoFromP2pAddr(m ma.Multiaddr) (*AddrInfo, error) {}

// AddrInfoFromString builds an AddrInfo from the string representation of a Multiaddr
func AddrInfoFromString(s string) (*AddrInfo, error) {}
