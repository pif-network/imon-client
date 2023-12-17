package core

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type AppError interface {
	Display() string
}

// Wrapper for internal errors, with concise cause.
type internalError struct {
	error error
	cause string
}

func InternalError(cause string, err error) *internalError {
	log.Error(fmt.Sprintf("[ERROR] Internal_Error: %s", cause))
	return &internalError{err, cause}
}
func (e *internalError) Error() string {
	return fmt.Sprintf("[ERROR] Internal_Error: %s", e.cause)
}
