package main

import (
	"net/http"

	"github.com/charmbracelet/log"
	"the-gorgeouses.com/imon-client/internal/core/server"
	"the-gorgeouses.com/imon-client/internal/modules"
)

func main() {
	appRouter := http.NewServeMux()
	modules.AttachRoutes(appRouter)

	log.Info("Starting server on port 8080")
	http.ListenAndServe(":8080", server.RemoveTrailingSlash(appRouter))
}
