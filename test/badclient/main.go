package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for i, a := range os.Args {
		fmt.Fprintf(os.Stderr, "printing %d\n", i)
		fmt.Fprintf(os.Stdout, "%s\n", a)
	}
	<-time.After(time.Minute)
}
