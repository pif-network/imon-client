package main

import (
	"log"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/modules/task"
	"the-gorgeouses.com/imon-client/internal/views"
)

func main() {
	http.HandleFunc("/", views.ServeRootView)
	http.HandleFunc("/post-key/", task.PostKeyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
