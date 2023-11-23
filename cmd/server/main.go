package main

import (
	"log"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/modules/task"
)

func main() {
	http.HandleFunc("/", task.ServeRootView)
	http.HandleFunc("/post-key/", task.PostKeyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
