package remora

import (
	"fmt"
)

type Errs []error

func (errs *Errs) Add(err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func (errs Errs) Result() error {
	if len(errs) > 0 {
		return &result{errs}
	}

	return nil
}

type result struct {
	errs Errs
}

func (res *result) Error() string {
	return fmt.Sprintf("%d errors", len(res.errs))
}
