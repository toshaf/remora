package binary

import (
	"github.com/toshaf/remora/comms"
	"os"
	"path"
	"syscall"
)

type provider struct {
	dir string
}

// Creates a new binary comms provider using the directory
// specified by dir to house the pipes
func NewProvider(dir string) comms.Provider {
	return &provider{
		dir: dir,
	}
}

// Creates the pipes and starts the server connections.
// Whether you get full or half duplex depends of the value of dir
func (p *provider) Server(name string, dir comms.Dir) (comms.In, comms.Out, error) {
	err := os.MkdirAll(p.dir, 0777)
	if err != nil {
		return nil, nil, err
	}

	base := path.Join(p.dir, name)

	var in comms.In
	if dir&comms.Input == comms.Input {
		// the file names are from the client's POV
		fname := base + ".out"
		os.Remove(fname)
		err = syscall.Mkfifo(fname, 0666)
		if err != nil {
			return nil, nil, err
		}

		in = NewServerIn(fname)
	}

	var out comms.Out
	if dir&comms.Output == comms.Output {
		// the file names are from the client's POV
		fname := base + ".in"
		os.Remove(fname)
		err = syscall.Mkfifo(fname, 0644)
		if err != nil {
			if in != nil {
				in.Close()
			}
			return nil, nil, err
		}
		out = NewServerOut(fname)
	}

	return in, out, nil
}

// Connects to the pipes created by the server.
// Whether you get full or half duplex depends of the value of dir
func (p *provider) Client(name string, dir comms.Dir) (comms.In, comms.Out, error) {
	base := path.Join(p.dir, name)

	var in comms.In
	if dir&comms.Input == comms.Input {
		// the file names are from the client's POV
		fname := base + ".in"

		in = NewClientIn(fname)
	}

	var out comms.Out
	if dir&comms.Output == comms.Output {
		// the file names are from the client's POV
		fname := base + ".out"

		out = NewClientOut(fname)
	}

	return in, out, nil
}
