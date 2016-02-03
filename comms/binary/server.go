package binary

import (
	"github.com/toshaf/remora/comms"
	"os"
)

type serverIn struct {
	name string
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the inbound half of a server connection.
func NewServerIn(name string) comms.In {
	reqs, errs, stop := makeRecv(name)
	return &serverIn{
		name: name,
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

func (s *serverIn) Recv(dest interface{}) error {
	return recv(s.reqs, dest)
}

func (s *serverIn) Close() error {
	s.stop.Close()
	<-s.errs
	return os.Remove(s.name)
}

type serverOut struct {
	name string
	reqs chan<- Request
	errs <-chan error
	stop Stopper
}

// Creates the outbound half of a server connection.
func NewServerOut(name string) comms.Out {
	reqs, errs, stop := makeSend(name)
	return &serverOut{
		name: name,
		reqs: reqs,
		errs: errs,
		stop: stop,
	}
}

func (s *serverOut) Send(src interface{}) error {
	return send(s.reqs, src)
}

func (s *serverOut) Close() error {
	s.stop.Close()
	<-s.errs
	return os.Remove(s.name)
}
