package crypto

import "github.com/libp2p/go-libp2p/core/crypto/pb"

// Key represents a crypto key that can be compared to another key
type Key interface {
	// Raw returns the raw bytes of the key (not wrapped in the
	// libp2p-crypto protobuf).
	//
	// This function is the inverse of {Priv,Pub}KeyUnmarshaler.
	Raw() ([]byte, error)

	// Type returns the protobuf key type.
	Type() pb.KeyType
}

// PubKey is a public key that can be used to verify data signed with the corresponding private key
type PubKey interface {
	Key
}

// PrivKey represents a private key that can be used to generate a public key and sign data
type PrivKey interface {
	Key

	// Return a public key paired with this private key
	GetPublic() PubKey
}
