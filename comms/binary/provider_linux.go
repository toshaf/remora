// +build linux
package binary

import (
	"os"
	"syscall"
)

func createFifo(fname string) error {
	os.Remove(fname)
	return syscall.Mkfifo(fname, 0666)
}
