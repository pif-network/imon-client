package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

type Task struct {
	Name        string
	Description string
	Completed   bool
}

func main() {
	serveRootView := func(w http.ResponseWriter, r *http.Request) {
		jsonStr := []byte(`{"key":"lily:0001"}`)
		res, err := http.Post("http://localhost:8000/v1/task-log", "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Printf("Failed to post to task-log service.")
			log.Printf(err.Error())
			tmpl := template.Must(template.ParseFiles("index.html"))
			tmpl.Execute(w, nil)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("Failed to read response body.")
			log.Printf(err.Error())
			tmpl := template.Must(template.ParseFiles("index.html"))
			tmpl.Execute(w, nil)
			return
		}

		log.Println(string(resBody))

		tmpl := template.Must(template.ParseFiles("index.html"))
		tasks := map[string][]Task{
			"Tasks": {
				{"Task 1", "Description 1", false},
				{"Task 2", "Description 2", false},
			},
		}

		tmpl.Execute(w, tasks)
	}

	addTaskHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("HTMX!!!")
		log.Println(r.PostFormValue("name"))
		log.Println(r.PostFormValue("description"))

		name := r.PostFormValue("name")
		description := r.PostFormValue("description")

		htmlStr := fmt.Sprintf(`
			<li>
				<div class="w-1/3 p-4 mb-2 border border-gray-400 rounded-lg">
					<h2>%s</h2>
					<p class="italic">%s</p>
				</div>
			</li>
		`, name, description)
		tmpl, _ := template.New("task").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", serveRootView)
	http.HandleFunc("/add-task/", addTaskHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
