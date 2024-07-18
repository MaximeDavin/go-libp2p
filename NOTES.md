## Goals

- simplify
- ownership

## Constraint

- inter-operability with all libp2p implementations
- same interface as go-libp2p; it can be used in prysm simply by switching a feature flag
- at least same performance

## libp2p features recap

| libp2p feature                | Use in Prysm | Notes                                                                                  |
| ----------------------------- | ------------ | -------------------------------------------------------------------------------------- |
| Ping                          | ❌           |                                                                                        |
| Addressing                    | ✔️           | use multiaddress                                                                       |
| TCP                           | ✔️           |                                                                                        |
| QUIC                          | ✔️           |                                                                                        |
| WebSocket                     | ❌           |                                                                                        |
| WebRTC                        | ❌           |                                                                                        |
| TLS                           | ✔️?          | Needed for quic ?                                                                      |
| Noise                         | ✔️           |                                                                                        |
| Early Multiplexer Negotiation | ✔️?          | [link](https://docs.libp2p.io/concepts/multiplex/early-negotiation/)                   |
| mplex                         | ✔️           |                                                                                        |
| yamux                         | ✔️           |                                                                                        |
| autonat                       | ✔️?          | not specified in Prysm, but enabled by default in go-libp2p                            |
| autonatv2                     | ❌           | not specified in Prysm, disabled by default in go-libp2p                               |
| circuit relay                 | ✔️?          | disabled by default in prysm                                                           |
| Hole Punching                 | ❌?          | not specified in Prysm, not sure if enabled by default in go-libp2p, I think it is not |
| UPNP                          | ✔️           | disabled by default in prysm, controled by NATPortMap option in go-libp2p              |
| boostrap discovery nodes      | ✔️           |                                                                                        |
| Kademlia                      | ❌           |                                                                                        |
| mdns                          | ❌           |                                                                                        |
| rendezvous                    | ❌           |                                                                                        |
| pubsub gossipsub              | ✔️           | no need to re-write ?                                                                  |

Maybe we should use [libp2p specs](https://github.com/libp2p/specs/tree/master), and note what we are going to keep, what we are going to remove, what we are going to add ?

## go-libp2p code tour

### Config

#### Removed unused options

go-libp2p offers a lot of options that we can remove because Prysm always use the same set of options.

Prysm construct a new libp2p object in two locations: `beacon-chain\p2p\options.go` and `cmd\prysmctl\p2p\client.go`.
Prysm use these options that can be configured:

- `libp2p.Identity`(set PeerKey)
- `ListenAddrs`
- `UserAgent`
- `ConnectionGater`
- `libp2p.Transport(libp2pquic.NewTransport)`
- `libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport) + Muxer(yamux.ID, yamux.DefaultTransport)`
- `libp2p.NATPortMap` (enable UPnP)
- `libp2p.AddrsFactory` (used to enable relay and to set explicit host addresses and dns)
- `libp2p.DisableRelay`
- `libp2p.ResourceManager` to disable libp2p default resource manager

Prysm always set these options to always have the same value:

- `libp2p.Transport(libp2ptcp.NewTCPTransport)` We always want to run TCP
- `libp2p.Security(noise.ID, noise.New)` We always use noise
- `libp2p.Ping(false)` We always disable Ping service

#### Simplify defaults

go-libp2p use `FallbackDefaults` to set defaults values if certain conditions are met. In my opinion we should remove this logic and use something simpler, like setting the default values explicitly in the constructor:

_TODO_ List options with default values that we are going to use

- `Peerstore := pstoremem.NewPeerstore()`

### New libp2p host

#### Fx

To instanciate a new host, go-libp2p use [Fx](https://uber-go.github.io/fx/) as a dependency injector (see `NewNode()` in `config\config.go`).
In my opinion, using Fx will help with dependency managment. The lifecycle hooks have a great potential too. BUT it makes the code way less explicit.
For reference [this is what the constructor was like before fx refactor](https://github.com/libp2p/go-libp2p/blob/c334288f8fe4d659f290043f788509e14f28cdde/config/config.go)

_TODO_ Write an example of the same code with and without Fx.

_TODO_ Decide if Fx is really needed or not.

### Components

Libp2p use various components / abstractions. We are going to list them and see if we should implement them or not

![components](https://camo.githubusercontent.com/8c3da2eadf623888f440368fbdb5a05c0a5a3c18ee10d5f2631dcb6482861c80/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f314676553747496d527362394776415744446f316c6538356a49726e464a4e56425f4f54505843313557774d2f7075623f683d343830)

#### Transport

The Transport interface allows to open connections to other peers by dialing them (implements Dial method), and also listen for incoming connections (implements Listen method).

We are implementing only two transports:

- tcp: go-libp2p use [this package](https://github.com/multiformats/go-multiaddr/blob/master/net/net.go) as a sort of wrapper of golang net package with multiaddress support.
- quic: _TODO_

#### Conn

A libp2p connection is a communication channel that allows peers to read and write data.

Connections between peers are established via transports, which can be thought of as “connection factories”. For example, the TCP transport allows you to create connections that use TCP/IP as their underlying substrate.

#### Muxer



#### Switch / Swarm

Used to be called Swarm.
A libp2p component responsible for composing multiple transports into a single interface, allowing application code to dial peers without having to specify what transport to use.
In addition to managing transports, the switch also coordinates the “connection upgrade” process, which promotes a “raw” connection from the transport layer into one that supports protocol negotiation, stream multiplexing, and secure communications.



#### Network

_TODO_

#### Peerstore

_TODO_

#### Event Bus

Allow communication between components. Components can publish and subscribe to specific events.

_TODO_ List events we need to keep

## Simplify

### Remove unused features

Libp2p implements A LOT of features that are not used in Prysm.

#### Example - Remove ping service

Libp2p ping service is always disabled when (https://github.com/MaximeDavin/prysm/blob/f2ce115ade55d2c01faa10d67f48cc8ca80e7bc3/beacon-chain/p2p/options.go)(Prysm instanciate a new libp2p object)

```go
	options := []libp2p.Option{
		...
		libp2p.Ping(false), // Ping Service is always disabled
	}
```

So we can remove the matching entry in the config and the [Ping service](https://github.com/libp2p/go-libp2p/blob/master/p2p/protocol/ping/ping.go) itself.

```go
type Config struct {
	...
	DisablePing bool // No need for this anymore
}
```

And finally the `p2p\protocol\ping` folder

### Rewrite and simplify

_TODO_
