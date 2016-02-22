// +build windows
package binary

import (
	"fmt"
)

func createFifo(name string) error {
	return fmt.Errorf("Not implemented for Windows")
}
