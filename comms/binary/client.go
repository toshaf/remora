package binary

import (
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/errors"
)

type client struct {
	in, out *half
}

func NewClient(files FNames) comms.Pipe {
	return &client{
		in:  newClientIn(files.In),
		out: newClientOut(files.Out),
	}
}

func (s *client) Close() error {
	s.in.stop.Close()
	s.out.stop.Close()

	errs := errors.Errors{}
	errs.Add(<-s.in.errs)
	errs.Add(<-s.out.errs)

	return errs.Result()
}

func (s *client) Send(src interface{}) error {
	return send(s.out.reqs, src)
}

func (s *client) Recv(dest interface{}) error {
	return recv(s.in.reqs, dest)
}

type half struct {
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the inbound half of a client connection.
func newClientIn(name string) *half {
	reqs, errs, stop := makeRecv(name)
	return &half{
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

// Creates the outbound half of a client connection.
func newClientOut(name string) *half {
	reqs, errs, stop := makeSend(name)
	return &half{
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}
