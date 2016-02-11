package remora

import (
	"github.com/toshaf/remora/errors"
)

type Closer interface {
	Close() error
}

type Closers []Closer

func (c Closers) CloseAll() errors.Errors {
	errs := errors.Errors{}
	for _, cl := range c {
		if cl != nil {
			errs.Add(cl.Close())
		}
	}

	return errs
}

func (c *Closers) Append(args ...Closer) {
	*c = append(*c, args...)
}
