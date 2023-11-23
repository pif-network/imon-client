package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	. "the-gorgeouses.com/imon-client/internal/core"
)

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

func GetUserTaskLogById(userKey string) (UserTaskLogResponse, error) {
	payload := fmt.Sprintf(`{"key": "%s"}`, userKey)
	res, err := http.Post(
		"http://localhost:8000/v1/task-log",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		log.Printf("Failed to post to task-log service.")
		log.Printf(err.Error())
		return UserTaskLogResponse{}, UpstreamError("Cannot reach service.")
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Failed to post to task-log service.")
		log.Printf("Status code: %d", res.StatusCode)
		return UserTaskLogResponse{}, UpstreamError("Invalid user key.")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return UserTaskLogResponse{}, InternalError("Failed to read response body.", err)
	}

	var resp UserTaskLogResponse
	if err := json.Unmarshal(resBody, &resp); err != nil {
		return UserTaskLogResponse{}, InternalError("Failed to unmarshal response body.", err)
	}

	return resp, nil
}
