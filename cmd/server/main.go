package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/core/server"
)

func main() {
	appRouter := http.NewServeMux()
	core.AttachRoutes(appRouter)

	log.Info("Starting server on port 8080")
	http.ListenAndServe(":8080", server.RemoveTrailingSlash(appRouter))
}
