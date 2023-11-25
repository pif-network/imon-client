package core

import (
	"net/http"

	"the-gorgeouses.com/imon-client/internal/modules/task"
	"the-gorgeouses.com/imon-client/internal/views"
)

func AttachRoutes(appRouter *http.ServeMux) {
	apiRouter := http.NewServeMux()
	apiRouter.Handle("/task/", http.StripPrefix("/task", task.GetTaskRouter()))

	appRouter.Handle("/api/", http.StripPrefix("/api", apiRouter))

	appRouter.HandleFunc("/", views.ServeRootView)
}
