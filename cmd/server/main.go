package main

import (
	"net/http"

	"the-gorgeouses.com/imon-client/internal/core"
)

func main() {
	appRouter := http.NewServeMux()
	core.AttachRoutes(appRouter)

	http.ListenAndServe(":8080", appRouter)
}
