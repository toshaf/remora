package remora

type Closer interface {
	Close() error
}

type Closers []Closer

func (c Closers) CloseAll() Errs {
	errs := Errs{}
	for _, cl := range c {
		errs.Add(cl.Close())
	}

	return errs
}

func (c *Closers) Append(args ...Closer) {
	*c = append(*c, args...)
}
