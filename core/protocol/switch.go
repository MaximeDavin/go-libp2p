package protocol

// Switch is the component responsible for "dispatching" incoming stream requests to
// their corresponding stream handlers. It is both a Negotiator and a Router.
type Switch interface {
	// Protocols returns a list of all registered protocol ID strings.
	// Note that the Router may be able to handle protocol IDs not
	// included in this list if handlers were added with match functions
	// using AddHandlerWithFunc.
	Protocols() []ID
}
