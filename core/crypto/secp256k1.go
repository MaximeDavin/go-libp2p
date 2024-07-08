package crypto

import (
	"io"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/libp2p/go-libp2p/core/crypto/pb"
)

// Secp256k1PrivateKey is a Secp256k1 private key
type Secp256k1PrivateKey secp256k1.PrivateKey

// Secp256k1PublicKey is a Secp256k1 public key
type Secp256k1PublicKey secp256k1.PublicKey

func GenerateSecp256k1Key(src io.Reader) (PrivKey, PubKey, error) {}

// Type returns the private key type
func (k *Secp256k1PrivateKey) Type() pb.KeyType {
	return pb.KeyType_Secp256k1
}

// UnmarshalSecp256k1PrivateKey returns a private key from bytes
func UnmarshalSecp256k1PrivateKey(data []byte) (k PrivKey, err error) {}

// Raw returns the bytes of the key
func (k *Secp256k1PrivateKey) Raw() ([]byte, error) {
	return (*secp256k1.PrivateKey)(k).Serialize(), nil
}

// GetPublic returns a public key
func (k *Secp256k1PrivateKey) GetPublic() PubKey {
	return (*Secp256k1PublicKey)((*secp256k1.PrivateKey)(k).PubKey())
}

// Type returns the public key type
func (k *Secp256k1PublicKey) Type() pb.KeyType {
	return pb.KeyType_Secp256k1
}

// Raw returns the bytes of the key
func (k *Secp256k1PublicKey) Raw() (res []byte, err error) {
}
