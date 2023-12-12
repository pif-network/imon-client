package server

import "fmt"

type UpstreamError struct {
	StatusCode int
	error
}

func NewUpstreamError(reason string, statusCode int) *UpstreamError {
	return &UpstreamError{
		StatusCode: statusCode,
		error:      fmt.Errorf(reason),
	}
}
func (e *UpstreamError) Error() string {
	return fmt.Sprintf("[ERROR] Upstream_Error: %s", e.error.Error())
}
func IsUpstreamError(err error) bool {
	_, ok := err.(*UpstreamError)
	return ok
}
func FixableByClient(err error) bool {
	if IsUpstreamError(err) {
		return err.(*UpstreamError).StatusCode == 400
	}
	return false
}
