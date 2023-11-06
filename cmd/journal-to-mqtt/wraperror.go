package main

import "fmt"

type WrappedError struct {
	Message string
	Err     error
}

func (w *WrappedError) Error() string {
	return fmt.Sprintf("%s: %v", w.Message, w.Err)
}
func WrapError(info string, err error) *WrappedError {
	return &WrappedError{
		Message: info,
		Err:     err,
	}
}
