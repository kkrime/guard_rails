package errors

import "fmt"

// this type of error gets reported to the user verbaitm
type RestError struct {
	Code int
	Err  string
}

func NewRestError(code int, errorMessage string, vars ...interface{}) error {
	return &RestError{
		Code: code,
		Err:  fmt.Sprintf(errorMessage, vars...),
	}
}

func (e *RestError) Error() string {
	return e.Err
}
