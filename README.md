# Prysm: Custom golang implementation of libp2p

https://hackmd.io/@6-HLeMXARN2tdFLKKcqrxw/rkU0eLmEC

## Goal

Empty minimal interface of go-libp2p used by Prysm. Every function or type used in Prysm has been copied from go-libp2p without its implementation.

## Usage

In Prysm go.mod use replace directive to link this version of go-libp2p ie

```
replace github.com/libp2p/go-libp2p => /path/to/custom/go-libp2p
```

## TODO

- Modify `core/crypto/pb/crypto.proto` (we do not need all KeyType) and generate a new `core/crypto/pb/crypto.pb.go`
