package server

import (
	"github.com/toshaf/remora"
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/comms/binary"
	"io"
	"os"
	"os/exec"
	"path"
)

type Args struct {
	Target string
}

type Run interface {
	Result() <-chan error
}

type run struct {
	result chan error
}

func (r run) Result() <-chan error {
	return r.result
}

type Server interface {
	Open(string) (comms.In, comms.Out, error)
	Start() (Run, error)
	Close() error
}

type server struct {
	provider comms.Provider
	target   string
	pipes    string
	closers  remora.Closers
}

func (s *server) Open(pipe string) (comms.In, comms.Out, error) {
	in, out, err := s.provider.Server(pipe, comms.Duplex)
	if err != nil {
		return nil, nil, err
	}
	s.closers.Append(in, out)

	return in, out, err
}

func New(args Args) Server {
	target := args.Target
	if !path.IsAbs(target) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		target = path.Join(wd, target)
	}
	pipes := path.Join(path.Dir(target), ".pipes")
	provider := binary.NewProvider(pipes)

	return &server{
		provider: provider,
		target:   target,
		pipes:    pipes,
	}
}

func (s *server) Start() (Run, error) {
	app := exec.Command(s.target, "-remora.pipes="+s.pipes)
	stdout, err := app.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := app.StderrPipe()
	if err != nil {
		return nil, err
	}
	err = app.Start()
	if err != nil {
		return nil, err
	}

	owait := copy(os.Stdout, stdout)
	ewait := copy(os.Stderr, stderr)

	result := make(chan error)

	go func() {
		<-owait
		<-ewait

		err := app.Wait()

		if err != nil {
			result <- err
		}
		close(result)
	}()

	return run{
		result: result,
	}, nil
}

func (s *server) Close() error {
	errs := s.closers.CloseAll()
	errs.Add(os.Remove(s.pipes))

	return errs.Result()
}

func copy(w io.Writer, r io.Reader) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		buffer := make([]byte, 1024)
		for {
			nr, err := r.Read(buffer)
			if nr == 0 {
				return
			}
			nw := 0
			for nw < nr {
				i, err := w.Write(buffer[nw:nr])
				if err != nil {
					panic(err)
				}
				nw += i
			}
			if err == io.EOF {
				return
			}
			if err != nil {
				panic(err)
			}
		}
	}()

	return wait
}
