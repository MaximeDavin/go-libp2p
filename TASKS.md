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
| ❌   |       | Main              | 3    | Use config to run and orchestrates services                                                    |
| ❌   |       | TCP transport     | 1    | Dial/Listen for TCP                                                                            |
| ❌   |       | quic transport    | 3?   | Dial/Listen for quic                                                                           |
| ❌   |       | mplex muxer       | 1    | stream multiplexer for mplex                                                                   |
| ❌   |       | yamux muxer       | 1    | stream multiplexer for yamux                                                                   |
| ❌   |       | Upgrader          | 1    | negotiate and upgrade connection: secured and multiplexed                                      |
| ❌   |       | ConnGater         | ?    | accept or reject connections when established or upgraded                                      |
| ❌   |       | ConnManager       | ?    | trims connections automatically                                                                |
| ❌   |       | Peerstore         | 2    | store addresses and other infos for every peerId                                               |
| ❌   |       | network/swarm     | 3    | open/close/store connections and streams                                                       |
| ❌   |       | host              | 3    | Main interface/ Highest level object                                                           |
| ❌   |       | Noise             | 4    | Maybe reuse ?                                                                                  |
| ❌   |       | Gossipsub         | 5    | Maybe reuse ?                                                                                  |
| ❌   |       | TLS               |      | Not sure if needed ? I think it is needed for quic                                             |
| ❌   |       | Metrics           |      | Not sure if needed ?                                                                           |
| ❌   |       | eventBus          |      | Not sure if needed ? Subscribe/Notify events between components Maybe use go-libp2p one as-is? |
| ❌   |       | Identify protocol |      | Not sure if needed ? Maybe use go-libp2p one as-is?                                            |

# Components

## Transport

```go
    type Transport interface{
        Dial(p peer.ID) Conn
        Listen() Listener
    }
```

## Network

```go
    type Network struct {
        conns       map[peer.ID][]*Conn
        listeners   map[Listener]struct{}
        transports  map[int]Transport
        muxers

    }

    func (n *netowrk) Connect(p peer.ID) Conn {
        if p in conns and conns[p].isUsable:
            return conns[p]
        else:
            conn = transport.Dial(p) // If we have quic + tcp which one should we use ?
            conns[p] = conn
            return conn
    }


```

## Host

```go
    type Host struct {
        network Network
        mux     Mux
        peerstore Peerstore
    }

    type HostI interface {

    }
```

# Workflow

## Build

1. Entrypoint `libp2p.New(opts ...Option)` is called in Prysm when a new Service is started. The returned object is stored in Service.Host. Options are store in a Config object.
2. `cfg.NewNode()` is called to build the host with the Config. Acts like a main function that build all services needed:

3. create a new Network/Swarm `n = newNetwork()`
4. create a new BasicHost `n.h = newHost(n)`
5. add transports to Network `for t in transport do n.addTransport(t)`. An upgrader is added to the transport object to add security (eg Noise) and stream multiplexing (mplex or yamux) when a new connection will be established.
6. start
