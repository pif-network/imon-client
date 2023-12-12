package server

import (
	"fmt"

	"the-gorgeouses.com/imon-client/internal/core"
)

type UpstreamError struct {
	statusCode int
	msg        string
	error
}

func NewUpstreamError(msg string, statusCode int, err error) *UpstreamError {
	return &UpstreamError{
		statusCode: statusCode,
		msg:        msg,
		error:      err,
	}
}
func (e *UpstreamError) Error() string {
	if e.error == nil {
		return fmt.Sprintf("Upstream_Error: %s", e.msg)
	} else {
		return fmt.Sprintf("Upstream_Error: %s", e.error.Error())
	}
}
func (e *UpstreamError) Msg() string {
	return e.msg
}
func IsUpstreamError(err error) bool {
	_, ok := err.(*UpstreamError)
	return ok
}
func FixableByClient(err error) (core.AppError, bool) {
	if uerr, ok := err.(*UpstreamError); ok {
		return uerr, uerr.statusCode == 400
	}
	return nil, false
}
