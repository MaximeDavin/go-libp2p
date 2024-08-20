# Tasks

| Symbol | State       |
| ------ | ----------- |
| ✔️     | Done        |
| 👷     | In progress |
| 👀     | PR open     |
| ❌     | Not started |

Size:
1(Trivial) -> 5(Very complex)

| Done | Owner | Component         | Size | Notes                                                                                          |
| ---- | ----- | ----------------- | ---- | ---------------------------------------------------------------------------------------------- |
| ❌   |       | Config            | 2    | Check and store options passed to libp2p.New()                                                 |
| ❌   |       | Main              | 2    | Use config to run and orchestrates services                                                    |
| ❌   |       | TCP transport     | 3    | Dial/Listen for TCP                                                                            |
| ❌   |       | quic transport    | 3?   | Dial/Listen for quic                                                                           |
| ❌   |       | TLS               | 2    | Needed for quic                                                                                |
| ❌   |       | mplex muxer       | 1    | stream multiplexer for mplex                                                                   |
| ❌   |       | yamux muxer       | 1    | stream multiplexer for yamux                                                                   |
| ❌   |       | Noise             | 4    | Maybe reuse ?                                                                                  |
| ❌   |       | Upgrader          | 1    | negotiate and upgrade connection: secured and multiplexed                                      |
| ❌   |       | ConnGater         | X    | accept or reject connections when established or upgraded (not needed ?)                       |
| ❌   |       | ConnManager       | X    | trims connections automatically (not needed ?)                                                 |
| ❌   |       | Peerstore         | 2    | store addresses and other infos for every peerId                                               |
| ❌   |       | network/swarm     | 3    | open/close/store connections and streams                                                       |
| ❌   |       | host              | 3    | Main interface/ Highest level object                                                           |
| ❌   |       | Gossipsub         | 5    | Maybe reuse ?                                                                                  |
| ❌   |       | Metrics           | X    | Not sure if needed ?                                                                           |
| ❌   |       | eventBus          | X    | Not sure if needed ? Subscribe/Notify events between components Maybe use go-libp2p one as-is? |
| ❌   |       | Identify protocol | X    | Not sure if needed ? Maybe use go-libp2p one as-is?                                            |

# Components

## Transport Layer

- Used to implements TCP and Quic.
- Open connections to other peer with Dial.
- Listen for incoming connections using a listener.
- Upgrade connections with security and multiplexing

```go
    type Transport interface {
        Dial(addr ma.Multiaddr) StreamConn
        Listen(addr ma.Multiaddr) Listener
    }

    type Listener interface {
        Accept() (StreamConn, error)
        Close() error
    }
```

### Connection TCP

1. Start with a 'raw' connection [that support multiaddrs](https://github.com/multiformats/go-multiaddr/blob/f63b0ed297ac01c14c4e428ef72bdabd7fcf9de2/net/net.go#L19): `manet.Conn`
2. Use noise to upgrade the `manet.Conn` to a `SecuredConn`.

```go
    type SecuredConn interface {
        manet.Conn
        LocalPeer() peer.ID
        RemotePeer() peer.ID
    }

    func SecureUpgrade(c manet.Conn) SecuredConn {
        // Noise stuff
    }
```

3. Use multiplexer to upgrade the `SecuredConn` to a `StreamConn` which add streaming capabilities.
   Use a `ListenStream()` method that will be run when a new StreamConn is established. It allows the peer to open an inbound stream over this connection. Similar to [swarm-conn.start method](https://github.com/libp2p/go-libp2p/blob/master/p2p/net/swarm/swarm_conn.go#L113)

```go
    type StreamConn interface {
        SecuredConn
        OpenStream() Stream
        AcceptStream() Stream
        Close()
        // Replace swarm_conn.start. When a new connection is created, it should listen for incoming stream. The stream is passed into the handler to handle a specific protocl
        ListenStream(handler func(Stream))
    }

    func MultiplexUpgrade(c SecuredConn) StreamConn {
        // Yamux stuff
    }
```

### Connection Quic

TODO

- Need to implement TLS
- Like TCP: it must implement the transport layer interface (dial, listen, accept, close) and return a `StreamConn`

### Noise

TODO
Used to upgrade the raw connection to a secure one. See `p2p/security/noise/` for more details.

## Network / Swarm

A connection manager responsible for opening/closing/retrieving connections.
It has a map that stores a connection for each peer.

```go
    type Network struct {
        conns       map[peer.ID]StreamConn // Maybe an array of connections for each peer ?
        listeners   []Listener
        mux     MultistreamMuxer // I think it makes more sense to have the muxer here instead of the host. For compatibility we can add a function in the host struct that return this mux (host.Mux() => network.mux)
        peerstore Peerstore // Same as mux
    }
    func DialPeer(ai peer.AddrInfo) StreamConn{
        // Check if there is a usable connection to this peer in conns map
        // true -> return the connection
        // false -> dial peer and add the connection to conns map then return it

        // absorb addresses into peerstore

        // Ideas:
        // - prioritize quic over tcp if supported
        // - add timeouts when dialing
        // - retry if it fails
        // - limit concurrent dials to avoid saturation
    }

    func Listen(addrs ...ma.Multiaddr) {
        // For each addr, a listener matching the associated transport of this addr is set up
        // The new listener is added to the listeners array
        // When the listener accepts a connection, it is added to the conns map
    }

    func NewStream(p peer.ID) Stream {
        // Call DialPeer() to get a connection to the peer
        // Call OpenStream on this connection
        // Maybe we should store the stream somewhere so it can be cleaned later

        // Negociate protocol with go-multistream.SelectOneOf
        // and store the protocol with Stream.SetProtocol()
    }

    func Close() {
        // Close and destroy conns map, listener array etc
    }

    // Set a stream handler to be used when a remote peer initiate a connection and starts a stream with the local peer.
    func SetStreamHandler(pid protocol.ID, handler func(Stream)) {
        // call mux.SetHandler(proto, handler)
    }


```

## Host

IMHO Host and Network could be merged together. If we decide to do so we can move as much logic as possible from Host to Network and keep Host for compatibility reasons, acting as a 'proxy' that call Network methods

```go
    type Host struct {
        network Network
    }

    func Connect(ai peer.AddrInfo) {
        // call network.DialPeer
    }

    func NewStream(p peer.ID, protos ...protocol.ID) Stream {
        // call network.NewStream
    }

    func Close() {
        // close and clean
    }
```

## Gossipsub

TODO: We need to decide if we use `go-libp2p-pubsub` package as is, or if we chose to re-implement it.
If you want to dig this subject feel free to update this part !

# Workflow

## Build

1. Entrypoint `libp2p.New(opts ...Option)` is called in Prysm when a new Service is started. The returned object is stored in Service.Host. Options are store in a Config object.
2. `cfg.NewNode()` is called to build the host with the Config. Acts like a main function that build all services needed:

3. create a new Network/Swarm `n = newNetwork()`
4. create a new BasicHost `n.h = newHost(n)`
5. add transports to Network `for t in transport do n.addTransport(t)`. An upgrader is added to the transport object to add security (eg Noise) and stream multiplexing (mplex or yamux) when a new connection will be established.
6. start
