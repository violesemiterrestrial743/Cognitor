package util

import "fmt"

type FieldError struct {
	Field string
	Err   error
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %v", e.Field, e.Err)
}

func (e FieldError) Unwrap() error {
	return e.Err
}
