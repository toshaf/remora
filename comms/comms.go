package comms

// Represents the outbound half of a duplex connection.
type Out interface {
	// Sends the value passed in, blocks until the send
	// completes and returns any error.
	// The value passed in could be a value or a pointer.
	Send(interface{}) error
	// Attempts to close this connection, returning any error.
	Close() error
}

// Represents the inbound half of a duplex connection.
type In interface {
	// Reads an incoming value from the pipe to the location
	// passed in; the value passed in must be a pointer.
	Recv(interface{}) error
	// Attempts to close this connection, returning any error.
	Close() error
}

// Represents a Remora comms provider
type Provider interface {
	// Creates the server's pipes in readiness for the client
	// to connect.
	// name is a logical name
	// dir specifies whether the connection should be in-only,
	// out-only or full duplex
	// if the error is non-nil, both returns will be nil
	Server(name string, dir Dir) (In, Out, error)
	// Connects to the pipes set up by the server
	// name is a logical name
	// dir specifies whether the connection should be in-only,
	// out-only or full duplex
	// if the error is non-nil, both returns will be nil
	Client(name string, dir Dir) (In, Out, error)
}
