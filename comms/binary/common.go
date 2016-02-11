package binary

import (
	"encoding/gob"
	"os"
)

type Request struct {
	target interface{}
	reply  chan error
}

type Stopper struct {
	ch chan<- struct{}
}

func (s *Stopper) Close() {
	if s.ch != nil {
		close(s.ch)
		s.ch = nil
	}
}

type FNames struct {
	In, Out string
}

func recv(ch chan<- Request, dest interface{}) error {
	reply := make(chan error)
	ch <- Request{dest, reply}
	return <-reply
}

func send(ch chan<- Request, src interface{}) error {
	reply := make(chan error)
	ch <- Request{src, reply}
	return <-reply
}

func makeTx(fname string, open func(string) (*os.File, error), tx func(interface{}) error) (chan<- Request, <-chan error, Stopper) {
	reqs := make(chan Request)
	errs := make(chan error)
	stop := make(chan struct{})

	go func() {
		var (
			err  error
			file *os.File
		)

		defer func() {
			if file != nil {
				errs <- file.Close()
			}
			close(errs)
		}()

		for {
			select {
			case <-stop:
				stop = nil
				return
			case req := <-reqs:
				if file == nil {
					file, err = open(fname)
					if err != nil {
						errs <- err
						return
					}
				}
				req.reply <- tx(req.target)
			}
		}
	}()

	return reqs, errs, Stopper{stop}
}

func makeRecv(fname string) (chan<- Request, <-chan error, Stopper) {
	var dec *gob.Decoder

	open := func(fname string) (*os.File, error) {
		file, err := os.Open(fname)
		if err == nil {
			dec = gob.NewDecoder(file)
		}

		return file, err
	}

	tx := func(v interface{}) error {
		return dec.Decode(v)
	}

	return makeTx(fname, open, tx)
}

func makeSend(fname string) (chan<- Request, <-chan error, Stopper) {
	var enc *gob.Encoder

	open := func(fname string) (*os.File, error) {
		file, err := os.OpenFile(fname, os.O_WRONLY, 0666)
		if err == nil {
			enc = gob.NewEncoder(file)
		}

		return file, err
	}

	tx := func(v interface{}) error {
		return enc.Encode(v)
	}

	return makeTx(fname, open, tx)
}
