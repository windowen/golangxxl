package apiresp

import (
	"errors"

	"liveJob/pkg/tools/errs"
)

type Checker interface {
	Check() error
}

func Validate(args interface{}) error {
	checker, ok := args.(Checker)
	if !ok {
		return errs.ErrArgs
	}

	if err := checker.Check(); err != nil {
		var codeErr errs.CodeError
		if errors.As(err, &codeErr) {
			return codeErr
		}
		return errs.ErrArgs.Wrap(err.Error())
	}

	return nil
}
