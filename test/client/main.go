package main

import (
	"fmt"
	"github.com/toshaf/remora/client"
	"github.com/toshaf/remora/test"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "creating client...")
	rc := client.New()
	fmt.Fprintf(os.Stderr, "done\n")
	defer func() {
		Check(rc.Close())
	}()

	fmt.Fprintf(os.Stderr, "connecting to maths ...")
	maths, err := rc.Connect("maths")
	Check(err)
	fmt.Fprintf(os.Stderr, "done\n")

	fmt.Fprintf(os.Stderr, "sending maths question...")
	err = maths.Send(test.Q{A: 5, B: 4, Op: "+"})
	Check(err)
	fmt.Fprintf(os.Stderr, "done\n")

	fmt.Fprintf(os.Stderr, "receiving answer...")
	var a test.A
	Check(maths.Recv(&a))
	fmt.Fprintf(os.Stderr, "done\n")

	fmt.Printf("Answer: %d\n", a.V)
	fmt.Fprintf(os.Stderr, "client finished\n")
}

func Check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
