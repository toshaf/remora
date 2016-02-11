package main

import (
	"fmt"
	"github.com/toshaf/remora/server"
	"os"
	"time"
)

func main() {
	srv := server.New(server.Args{Target: os.Args[1]})
	defer func() {
		Check(srv.Close())
	}()

	run, err := srv.Start("one", "two", "three")
	Check(err)

	select {
	case err := <-run.Result():
		Check(err)
	case <-time.After(time.Second):
		Check(run.Kill())
	}

	fmt.Println("Done")
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
