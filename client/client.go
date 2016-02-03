package client

import (
	"flag"
	"github.com/toshaf/remora"
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/comms/binary"
)

// Exposes client capabilities to a Remora client app
type Client interface {
	// Given a logical name, creates a duplex connection to the server.
	// The server must have created the relevant pipes also.
	// If the error is not nil, both halves of the pipe will be.
	// Either half of the pipe can be closed independently of the other, if it's
	// not required by the client/server protocol.
	Connect(pipe string) (comms.In, comms.Out, error)
	// Closes all connections maintained by this instance.
	Close() error
}

var pipes *string

func init() {
	pipes = flag.String("remora.pipes", ".", "The directory in which the host will set up named pipes")

	flag.Parse()
}

// Creates a new Client instance using the directory specified by the server
// process that started this process.
func New() Client {
	return &client{
		Provider: binary.NewProvider(*pipes),
	}
}

type client struct {
	Provider comms.Provider
	closers  remora.Closers
}

func (c *client) Connect(pipe string) (comms.In, comms.Out, error) {
	in, out, err := c.Provider.Client(pipe, comms.Duplex)
	if err != nil {
		return in, out, err
	}
	c.closers.Append(in, out)

	return in, out, err
}

func (c *client) Close() error {
	return c.closers.CloseAll().Result()
}
