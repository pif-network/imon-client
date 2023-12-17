package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		return UserTaskLogResponse{}, server.NewUpstreamError("Cannot reach service", http.StatusInternalServerError, err)
	}
	if res.StatusCode != http.StatusOK {
		// NOTE: The only not-ok status that this client is currently able to cause is 400.
		logger.Debug("upstream_response", "code", res.StatusCode)
		if bBody, err := io.ReadAll(res.Body); err != nil {
			logger.Error(err.Error())
		} else {
			logger.Debug("upstream_response", "body", string(bBody))
		}
		return UserTaskLogResponse{}, server.NewUpstreamError(
			"[Upstream_Error] Invalid user key.", http.StatusBadRequest, fmt.Errorf("Invalid user key"),
		)
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
		return AllUserRecordsResponse{}, server.NewUpstreamError("Cannot reach service", http.StatusInternalServerError, err)
	}
	if res.StatusCode != http.StatusOK {
		return AllUserRecordsResponse{}, server.NewUpstreamError("Invalid user key.", http.StatusBadRequest, err)
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
		return server.NewUpstreamError("Cannot reach service", http.StatusInternalServerError, err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return server.NewUpstreamError("Invalid user key.", http.StatusBadRequest, err)
	}
	defer res.Body.Close()

	return nil
}
