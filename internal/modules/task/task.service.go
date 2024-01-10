package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/core"
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
	resp, err := http.Post(
		"http://localhost:8000/v1/record",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return handleErrorHttpResponse[UserTaskLogResponse](resp, err)
	}
	if resp.StatusCode != http.StatusOK {
		return handleNotOkHttpResponse[UserTaskLogResponse](resp)
	}
	defer resp.Body.Close()

	var dResp UserTaskLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&dResp); err != nil {
		logger.Debug(err.Error())
		return UserTaskLogResponse{}, core.NewInternalError("Failed to unmarshal response body", err)
	}

	return dResp, nil
}

type AllUserRecordsResponse struct {
	Data struct {
		UserRecords []TaskLog `json:"user_records"`
	} `json:"data"`
	Status string `json:"status"`
}

func GetAllUserRecords() (AllUserRecordsResponse, error) {
	resp, err := http.Get("http://localhost:8000/v1/record/all")
	if err != nil {
		return handleErrorHttpResponse[AllUserRecordsResponse](resp, err)
	}
	if resp.StatusCode != http.StatusOK {
		return handleNotOkHttpResponse[AllUserRecordsResponse](resp)
	}
	defer resp.Body.Close()

	var dResp AllUserRecordsResponse
	if err := json.NewDecoder(resp.Body).Decode(&dResp); err != nil {
		logger.Debug(err.Error())
		return AllUserRecordsResponse{}, core.NewInternalError("Failed to unmarshal response body", err)
	}

	return dResp, nil
}

func UpdateCurrentTask(userKey string, taskState TaskState) error {
	payload := fmt.Sprintf(`{"key": "%s", "state": "%s"}`, userKey, taskState)
	logger.Info(payload)

	resp, err := http.Post(
		"http://localhost:8000/v1/task/update",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		_, err = handleErrorHttpResponse[interface{}](resp, err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		_, err := handleNotOkHttpResponse[interface{}](resp)
		return err
	}
	defer resp.Body.Close()

	return nil
}
