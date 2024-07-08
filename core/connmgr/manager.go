package connmgr

// ConnManager tracks connections to peers, and allows consumers to associate
// metadata with each peer.
//
// It enables connections to be trimmed based on implementation-defined
// heuristics. The ConnManager allows libp2p to enforce an upper bound on the
// total number of open connections.
//
// ConnManagers supporting decaying tags implement Decayer. Use the
// SupportsDecay function to safely cast an instance to Decayer, if supported.
type ConnManager interface{}
