package errs

import "errors"

func RpcErrCheck(err error, checkErr CodeError) bool {
	unwrap := Unwrap(err)
	var codeErr CodeError
	if errors.As(unwrap, &codeErr) {
		if checkErr.Code() == codeErr.Code() {
			return true
		}
	}

	return false
}
