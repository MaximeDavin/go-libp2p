package transport

// Upgrader is a multistream upgrader that can upgrade an underlying connection
// to a full transport connection (secure and multiplexed).
type Upgrader interface{}

// Transport represents any device by which you can connect to and accept
// connections from other peers.
//
// The Transport interface allows you to open connections to other peers
// by dialing them, and also lets you listen for incoming connections.
//
// Connections returned by Dial and passed into Listeners are of type
// CapableConn, which means that they have been upgraded to support
// stream multiplexing and connection security (encryption and authentication).
//
// If a transport implements `io.Closer` (optional), libp2p will call `Close` on
// shutdown. NOTE: `Dial` and `Listen` may be called after or concurrently with
// `Close`.
//
// For a conceptual overview, see https://docs.libp2p.io/concepts/transport/
type Transport interface{}
