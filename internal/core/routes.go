package core

import (
	"net/http"

	"the-gorgeouses.com/imon-client/internal/modules/task"
	"the-gorgeouses.com/imon-client/internal/views"
)

func AttachRoutes(appRouter *http.ServeMux) {
	apiRouter := http.NewServeMux()

	apiRouter.Handle("/api/task/", task.GetTaskRouter())

	appRouter.Handle("/api/", apiRouter)
	appRouter.HandleFunc("/", views.ServeRootView)
}
