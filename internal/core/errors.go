package errors

import "fmt"

type upstreamError struct {
	error
}

func UpstreamError(reason string) *upstreamError {
	return &upstreamError{fmt.Errorf(reason)}
}
func (e *upstreamError) Error() string {
	return fmt.Sprintf("[ERROR] Upstream_Error: %s", e.error.Error())
}
func IsUpstreamError(err error) bool {
	_, ok := err.(*upstreamError)
	return ok
}

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
