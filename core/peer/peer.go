package peer

import ic "github.com/libp2p/go-libp2p/core/crypto"

// ID is a libp2p peer identity.
//
// Peer IDs are derived by hashing a peer's public key and encoding the
// hash output as a multihash. See IDFromPublicKey for details.
type ID string

// IDSlice for sorting peers.
type IDSlice []ID

func (id ID) String() string {}

// IDFromPublicKey returns the Peer ID corresponding to the public key pk.
func IDFromPublicKey(pk ic.PubKey) (ID, error) {}

// IDFromBytes casts a byte slice to the ID type, and validates
// the value to make sure it is a multihash.
func IDFromBytes(b []byte) (ID, error) {}

// Decode accepts an encoded peer ID and returns the decoded ID if the input is
// valid.
//
// The encoded peer ID can either be a CID of a key or a raw multihash (identity
// or sha256-256).
func Decode(s string) (ID, error) {}
