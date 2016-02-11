package binary

import (
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/errors"
	"os"
)

type server struct {
	in  *sHalf
	out *sHalf
}

func NewServer(files FNames) comms.Pipe {
	return &server{
		in:  newServerIn(files.In),
		out: newServerOut(files.Out),
	}
}

type sHalf struct {
	name string
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the inbound half of a server connection.
func newServerIn(name string) *sHalf {
	reqs, errs, stop := makeRecv(name)
	return &sHalf{
		name: name,
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

// Creates the outbound half of a server connection.
func newServerOut(name string) *sHalf {
	reqs, errs, stop := makeSend(name)
	return &sHalf{
		name: name,
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

func (s *server) Recv(dest interface{}) error {
	return recv(s.in.reqs, dest)
}

func (s *server) Send(src interface{}) error {
	return send(s.out.reqs, src)
}

func (s *server) Close() error {
	s.in.stop.Close()
	s.out.stop.Close()

	errs := errors.Errors{}
	errs.Add(<-s.out.errs)
	errs.Add(<-s.in.errs)

	if err := os.Remove(s.in.name); !os.IsNotExist(err) {
		errs.Add(err)
	}
	if err := os.Remove(s.out.name); !os.IsNotExist(err) {
		errs.Add(err)
	}

	return errs.Result()
}
