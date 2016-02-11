package main

import "time"

func main() {
	<-time.After(time.Minute)
}
