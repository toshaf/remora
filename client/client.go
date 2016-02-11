package client

import (
	"github.com/toshaf/remora"
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/comms/binary"
	"os"
	"strings"
)

var pipes string

func init() {
	args := []string{}
	for _, a := range os.Args {
		if strings.HasPrefix(a, "--remora.pipes=") {
			pipes = strings.Split(a, "=")[1]
		} else {
			args = append(args, a)
		}
	}
}

// Creates a new Client instance using the directory specified by the server
// process that started this process.
func New() remora.Client {
	return &client{
		Provider: binary.NewProvider(pipes),
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
