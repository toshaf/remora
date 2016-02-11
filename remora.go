package remora

import (
	"github.com/toshaf/remora/comms"
)

// Exposes client capabilities to a Remora client app
type Client interface {
	// Given a logical name, creates a duplex connection to the server.
	// The server must have created the relevant pipes also.
	// If the error is not nil, both halves of the pipe will be.
	// Either half of the pipe can be closed independently of the other, if it's
	// not required by the client/server protocol.
	Connect(pipe string) (comms.Pipe, error)
	// Closes all connections maintained by this instance.
	Close() error
}
