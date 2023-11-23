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
