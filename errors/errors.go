package errors

import (
	"bytes"
	"fmt"
)

type Errors []error

func (errs *Errors) Add(err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func (errs Errors) Result() error {
	if len(errs) > 0 {
		return &result{errs}
	}

	return nil
}

type result struct {
	errs Errors
}

func (res *result) Error() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "%d errors", len(res.errs))
	for _, err := range res.errs {
		fmt.Fprintf(&buf, err.Error())
	}
	return buf.String()
}
