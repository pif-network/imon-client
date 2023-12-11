package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/core/server"
)

type TaskState string

const (
	Begin TaskState = "Begin"
	Break TaskState = "Break"
	Back  TaskState = "Back"
	End   TaskState = "End"
	Idle  TaskState = "Idle"
)

func (t TaskState) String() string {
	return string(t)
}

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
	UserName    string `json:"user_name"`
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
		return UserTaskLogResponse{}, server.UpstreamError("Cannot reach service.")
	}
	if res.StatusCode != http.StatusOK {
		return UserTaskLogResponse{}, server.UpstreamError("Invalid user key.")
	}
	defer res.Body.Close()

	var resp UserTaskLogResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		fmt.Println("error", err)
		return UserTaskLogResponse{}, core.InternalError("Failed to unmarshal response body.", err)
	}

	return resp, nil
}

type AllUserRecordsResponse struct {
	Data struct {
		UserRecords []TaskLog `json:"user_records"`
	} `json:"data"`
	Status string `json:"status"`
}

func GetAllUserRecords() (AllUserRecordsResponse, error) {
	res, err := http.Get("http://localhost:8000/v1/record/all")
	if err != nil {
		return AllUserRecordsResponse{}, server.UpstreamError("Cannot reach service.")
	}
	if res.StatusCode != http.StatusOK {
		return AllUserRecordsResponse{}, server.UpstreamError("Invalid user key.")
	}
	defer res.Body.Close()

	var resp AllUserRecordsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		fmt.Println("error", err)
		return AllUserRecordsResponse{}, core.InternalError("Failed to unmarshal response body.", err)
	}

	return resp, nil
}

func UpdateCurrentTask(userKey string, taskState TaskState) error {
	payload := fmt.Sprintf(`{"key": "%s", "state": "%s"}`, userKey, taskState)
	fmt.Println(payload)
	res, err := http.Post(
		"http://localhost:8000/v1/task/update",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return server.UpstreamError("Cannot reach service.")
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return server.UpstreamError("Invalid user key.")
	}
	defer res.Body.Close()

	return nil
}
