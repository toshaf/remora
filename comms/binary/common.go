package binary

import (
	"encoding/gob"
	"os"
)

type Request struct {
	target interface{}
	reply  chan error
}

func recv(ch chan<- Request, dest interface{}) error {
	reply := make(chan error)
	ch <- Request{dest, reply}
	return <-reply
}

// Used
type Stopper struct {
	ch chan<- struct{}
}

func (s Stopper) Close() {
	if s.ch != nil {
		close(s.ch)
		s.ch = nil
	}
}

func makeRecv(fname string) (chan<- Request, <-chan error, Stopper) {
	reqs := make(chan Request)
	errs := make(chan error)
	stop := make(chan struct{})

	go func() {
		file, err := os.Open(fname)
		if err != nil {
			errs <- err
			return
		}
		defer file.Close()
		dec := gob.NewDecoder(file)

		for {
			select {
			case <-stop:
				close(errs)
				return
			case req := <-reqs:
				req.reply <- dec.Decode(req.target)
			}
		}
	}()

	return reqs, errs, Stopper{stop}
}

func send(ch chan<- Request, src interface{}) error {
	reply := make(chan error)
	ch <- Request{src, reply}
	return <-reply
}

func makeSend(fname string) (chan<- Request, <-chan error, Stopper) {
	reqs := make(chan Request)
	errs := make(chan error)
	stop := make(chan struct{})

	go func() {
		file, err := os.OpenFile(fname, os.O_WRONLY, 0666)
		if err != nil {
			errs <- err
			return
		}
		defer file.Close()
		enc := gob.NewEncoder(file)

		for {
			select {
			case <-stop:
				close(errs)
				return
			case req := <-reqs:
				req.reply <- enc.Encode(req.target)
			}
		}
	}()

	return reqs, errs, Stopper{stop}
}
