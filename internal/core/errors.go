package core

import (
	"fmt"
)

type AppError interface {
	// Display returns a user-friendly message.
	Display() string
}

// Wrapper for internal errors, with concise cause.
type internalError struct {
	msg string
	error
}

func NewInternalError(msg string, err error) *internalError {
	return &internalError{
		msg:   msg,
		error: err,
	}
}
func (e *internalError) Error() string {
	if e.error == nil {
		return fmt.Sprintf("Internal_Error: %s", e.msg)
	} else {
		return fmt.Sprintf("Internal_Error: %s", e.error.Error())
	}
}
func (e *internalError) Display() string {
	return fmt.Sprintf("[Internal_Error] %s.", e.msg)
}
