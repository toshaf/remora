package server

import (
	"fmt"
	"github.com/toshaf/remora"
	"github.com/toshaf/remora/comms"
	"github.com/toshaf/remora/comms/binary"
	"io"
	"os"
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
	Open(string) (comms.Pipe, error)
	Start() (Run, error)
	Close() error
}

type server struct {
	provider comms.Provider
	target   string
	pipes    string
	closers  remora.Closers
}

func (s *server) Open(name string) (comms.Pipe, error) {
	pipe, err := s.provider.Server(name)
	if err != nil {
		return nil, err
	}
	s.closers.Append(pipe)

	return pipe, err
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
	attr := &os.ProcAttr{
		Dir:   path.Dir(s.target),                         // start the process in the binary's directory
		Env:   nil,                                        // use the environment
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, // pass the output through
		Sys:   nil,                                        // no special requirements yet
	}
	args := []string{
		"--remora.pipes=" + s.pipes,
	}
	app, err := os.StartProcess(s.target, args, attr)
	if err != nil {
		return nil, err
	}

	result := make(chan error)

	go func() {
		state, err := app.Wait()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't wait on process: %s\n", err)
			err = app.Kill()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't kill process: %s\n", err)
			}
			err = app.Release()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't release process: %s\n", err)
			}
			panic(fmt.Errorf("Lost control of process"))
		}

		if !state.Success() {
			result <- fmt.Errorf(state.String())
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
