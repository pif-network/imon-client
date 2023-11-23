package main

import (
	"log"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ServeRootView)
	http.HandleFunc("/post-key/", handlers.PostKeyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
