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

// Creates the pipe and starts the server connections.
func (p *provider) Server(name string) (comms.Pipe, error) {
	err := os.MkdirAll(p.dir, 0777)
	if err != nil {
		return nil, err
	}

	base := path.Join(p.dir, name)

	// the file names are from the client's POV
	iname := base + ".out"
	err = createFifo(iname)
	if err != nil {
		return nil, err
	}

	oname := base + ".in"
	err = createFifo(oname)
	if err != nil {
		return nil, err
	}

	return NewServer(FNames{In: iname, Out: oname}), nil
}

// Connects to the pipes created by the server.
func (p *provider) Client(name string) (comms.Pipe, error) {
	base := path.Join(p.dir, name)

	// the file names are from the client's POV
	iname := base + ".in"
	_, err := os.Stat(iname)
	if err != nil {
		return nil, err
	}
	oname := base + ".out"
	_, err = os.Stat(oname)
	if err != nil {
		return nil, err
	}

	return NewClient(FNames{In: iname, Out: oname}), nil
}

func createFifo(fname string) error {
	os.Remove(fname)
	return syscall.Mkfifo(fname, 0666)
}
