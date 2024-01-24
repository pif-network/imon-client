package task

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/core/server"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "task_module",
	Level:           log.DebugLevel,
})

var UpstreamEndpoint = struct {
	User string
	Sudo string
}{
	User: "http://localhost:8000/v1/rpc/user",
	Sudo: "http://localhost:8000/v1/rpc/sudo",
}

func handleNotOkHttpResponse[T interface{}](res *http.Response) (T, error) {
	// NOTE: The only not-ok status that this client is currently able to cause is 400.
	logger.Debug("upstream_response", "code", res.StatusCode)

	var t T
	if bBody, err := io.ReadAll(res.Body); err != nil {
		logger.Debug(err.Error())
		return t, core.NewInternalError(
			"Cannot read request body", err,
		)
	} else {
		logger.Debug("upstream_response", "body", string(bBody))
	}
	return t, server.NewUpstreamError(
		"Invalid user key", http.StatusBadRequest, nil,
	)
}

func handleErrorHttpResponse[T interface{}](res *http.Response, error error) (T, error) {
	var t T
	return t, core.NewInternalError("Failed to perform HTTP request.", error)
}
