package binary

import (
	"github.com/toshaf/remora/comms"
)

type clientIn struct {
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the inbound half of a client connection.
func NewClientIn(name string) comms.In {
	reqs, errs, stop := makeRecv(name)
	return &clientIn{
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

func (s *clientIn) Recv(dest interface{}) error {
	return recv(s.reqs, dest)
}

func (s *clientIn) Close() error {
	s.stop.Close()
	return <-s.errs
}

type clientOut struct {
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the outbound half of a client connection.
func NewClientOut(name string) comms.Out {
	reqs, errs, stop := makeSend(name)
	return &clientOut{
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

func (s *clientOut) Send(src interface{}) error {
	return send(s.reqs, src)
}

func (s *clientOut) Close() error {
	s.stop.Close()
	return <-s.errs
}
