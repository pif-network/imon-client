package main

import (
	"net/http"

	"the-gorgeouses.com/imon-client/internal/modules/task"
	"the-gorgeouses.com/imon-client/internal/views"
)

func main() {
	apiRouter := http.NewServeMux()
	appRouter := http.NewServeMux()

	appRouter.HandleFunc("/", views.ServeRootView)

	apiRouter.HandleFunc("/api/task/post-key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			task.PostKeyHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	appRouter.Handle("/api/", apiRouter)

	http.ListenAndServe(":8080", appRouter)
}
