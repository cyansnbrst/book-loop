package validator

import "errors"

var (
	ErrJSONIsNotValid = errors.New("validation error")
)

type ValidationError struct {
	Errors map[string]string
	Err    error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}
