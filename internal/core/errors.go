package core

import "fmt"

// Wrapper for internal errors, with concise cause.
type internalError struct {
	error error
	cause string
}

func InternalError(cause string, err error) *internalError {
	return &internalError{err, cause}
}
func (e *internalError) Error() string {
	return fmt.Sprintf("[ERROR] Internal_Error: %s", e.cause)
}
