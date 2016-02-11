package client

import (
	"flag"
	"github.com/toshaf/remora"
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/comms/binary"
)

var pipes *string

func init() {
	pipes = flag.String("remora.pipes", ".", "The directory in which the host will set up named pipes")

	flag.Parse()
}

// Creates a new Client instance using the directory specified by the server
// process that started this process.
func New() remora.Client {
	return &client{
		Provider: binary.NewProvider(*pipes),
	}
}

type client struct {
	Provider comms.Provider
	closers  remora.Closers
}

func (c *client) Connect(name string) (comms.Pipe, error) {
	pipe, err := c.Provider.Client(name)
	if err != nil {
		return pipe, err
	}
	c.closers.Append(pipe)

	return pipe, err
}

func (c *client) Close() error {
	return c.closers.CloseAll().Result()
}
