package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/views/pages"
)

func ServeRootView(w http.ResponseWriter, r *http.Request) {
	_ = pages.Index().Render(r.Context(), w)
}

type TaskState string

const (
	Begin TaskState = "Begin"
	Break TaskState = "Break"
	Back  TaskState = "Back"
	End   TaskState = "End"
	Idle  TaskState = "Idle"
)

type Task struct {
	BeginTime string    `json:"begin_time"`
	Duration  int       `json:"duration"`
	EndTime   string    `json:"end_time"`
	Name      string    `json:"name"`
	State     TaskState `json:"state"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type TaskLog struct {
	ID          int    `json:"id"`
	CurrentTask Task   `json:"current_task"`
	TaskHistory []Task `json:"task_history"`
}

func (t *TaskLog) UnmarshalJSON(data []byte) error {
	type Alias TaskLog
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type UserTaskLogResponse struct {
	Data struct {
		TaskLog TaskLog `json:"task_log"`
	} `json:"data"`
	Status string `json:"status"`
}

func PostKeyHandler(w http.ResponseWriter, r *http.Request) {
	userKey := r.PostFormValue("user-key")
	log.Printf("User key: %s", userKey)

	payload := fmt.Sprintf(`{"key": "%s"}`, userKey)
	res, err := http.Post(
		"http://localhost:8000/v1/task-log",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		log.Printf("Failed to post to task-log service.")
		log.Printf(err.Error())
		// templates.ExecuteTemplate(w, "invalid-user-key", "Cannot reach service.")
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Failed to post to task-log service.")
		log.Printf("Status code: %d", res.StatusCode)
		// templates.ExecuteTemplate(w, "invalid-user-key", "Invalid user key.")
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read response body.")
		log.Printf(err.Error())
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	var resp UserTaskLogResponse
	if err := json.Unmarshal(resBody, &resp); err != nil {
		log.Printf("Failed to unmarshal response body.")
		log.Printf(err.Error())
	}

	// if err := templates.ExecuteTemplate(w, "task-list", resp); err != nil {
	// 	log.Printf("Failed to execute template.")
	// 	log.Printf(err.Error())
	// }
}
