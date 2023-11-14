package main

import (
	"fmt"
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
	fmt.Println("Hello, World!")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tasks := map[string][]Task{
			"Tasks": {
				{"Task 1", "Description 1", false},
				{"Task 2", "Description 2", false},
			},
		}

		tmpl.Execute(w, tasks)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-task/", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
