package main

import (
	"fmt"
	"github.com/toshaf/remora/server"
	"github.com/toshaf/remora/test"
	"io"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <app>", args[0])
		os.Exit(1)
	}

	srv := server.New(server.Args{Target: args[1]})
	defer func() {
		Check(srv.Close())
	}()

	maths(srv)

	run, err := srv.Start()
	Check(err)

	err = <-run.Result()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func maths(srv server.Server) {
	in, out, err := srv.Open("maths")
	Check(err)

	go func() {
		for {
			q := test.Q{}
			err := in.Recv(&q)
			if err == io.EOF {
				return
			}
			Check(err)
			Check(out.Send(runQ(q)))
		}
	}()
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func runQ(q test.Q) test.A {
	var a test.A
	switch q.Op {
	case "+":
		a.V = q.A + q.B
	case "-":
		a.V = q.A - q.B
	case "*":
		a.V = q.A * q.B
	case "/":
		a.V = q.A / q.B
	default:
		a.V = 42
	}

	return a
}
