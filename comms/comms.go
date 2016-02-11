package comms

// Represents a full duplex connection.
type Pipe interface {
	// Sends the value passed in, blocks until the send
	// completes and returns any error.
	// The value passed in could be a value or a pointer.
	Send(interface{}) error
	// Reads an incoming value from the pipe to the location
	// passed in; the value passed in must be a pointer.
	Recv(interface{}) error
	// Attempts to close the pipe, returns an error if it can't
	Close() error
}

// Represents a Remora comms provider
type Provider interface {
	// Creates the server's pipe in readiness for the client
	// to connect.
	// name is a logical name
	// if the error is non-nil, the Pipe will be nil
	Server(name string) (Pipe, error)
	// Connects to the pipe set up by the server
	// name is a logical name
	// if the error is non-nil, the Pipe will be nil
	Client(name string) (Pipe, error)
}
