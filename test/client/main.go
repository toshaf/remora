package main

import (
	"fmt"
	"github.com/toshaf/remora/client"
	"github.com/toshaf/remora/test"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "client\n")

	rc := client.New()
	defer func() {
		Check(rc.Close())
	}()

	in, out, err := rc.Connect("maths")
	Check(err)

	err = out.Send(test.Q{A: 5, B: 4, Op: "+"})
	Check(err)

	var a test.A
	Check(in.Recv(&a))

	fmt.Printf("Answer: %d\n", a.V)
	fmt.Fprintf(os.Stderr, "client finished\n")
}

func Check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
