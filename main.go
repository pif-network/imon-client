package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

type TaskState int

const (
	Begin TaskState = iota
	Break
	Back
	End
	Idle
)

type Task struct {
	Name      string
	TaskState TaskState
	BeginTime int64
	EndTime   int64
	Duration  int64
}

type TaskLog struct {
	CurrentTask Task
	Id          string
	TaskHistory []Task
}

type UserTaskLogResponse struct {
	Status string
	Data   TaskLog
}

func main() {
	serveRootView := func(w http.ResponseWriter, r *http.Request) {
		jsonStr := []byte(`{"key":"lily:0001"}`)
		res, err := http.Post(
			"http://localhost:8000/v1/task-log",
			"application/json",
			bytes.NewBuffer(jsonStr),
		)
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
		resp := UserTaskLogResponse{}
		json.Unmarshal(resBody, &resp)

		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.Execute(w, resp.Data.CurrentTask)
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
