package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/core/server"
)

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

type User struct {
	userKey  string
	userType string
	name     string
	id       int
}

type RpcPayload struct {
	Metadata struct {
		Of string `json:"of"`
	} `json:"metadata"`
	Payload map[string]interface{} `json:"payload"`
}

func generatePayload(userKey string, eventType string) string {
	payload := RpcPayload{
		Metadata: struct {
			Of string `json:"of"`
		}{
			Of: "sudo",
		},
		Payload: map[string]interface{}{
			"key":        userKey,
			"event_type": eventType,
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	return string(b)
}

type SingleRecordResponse struct {
	Data struct {
		Id             int     `json:"id"`
		Name           string  `json:"name"`
		PublishedTasks []STask `json:"published_tasks"`
	}
	Status string `json:"status"`
}

func GetSingleRecordSudo(userKey string) (SingleRecordResponse, error) {
	payload := generatePayload(userKey, UpstreamEventType.GetSingleRecord)
	logger.Debug(payload)
	resp, err := http.Post(
		"http://localhost:8000/v1/rpc/sudo",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return handleErrorHttpResponse[SingleRecordResponse](resp, err)
	}
	if resp.StatusCode != http.StatusOK {
		return handleNotOkHttpResponse[SingleRecordResponse](resp)
	}
	defer resp.Body.Close()

	var dResp SingleRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&dResp); err != nil {
		logger.Debug(err.Error())
		return SingleRecordResponse{}, core.NewInternalError("Failed to unmarshal response body", err)
	}

	return dResp, nil
}

// RefreshData refreshes the data of the record, then renders them accordingly.
func (u User) RefreshData(w http.ResponseWriter, r *http.Request) error {
	switch u.userType {
	case "user":
		respTaskLog, err := GetUserTaskLogById(u.userKey)
		if err != nil {
			return err
		}
		_ = CurrentTaskAndExecutionLog(respTaskLog.Data.TaskLog).Render(r.Context(), w)
		respAllRecords, err := GetAllUserRecords()
		if err != nil {
			return err
		}
		_ = ActiveUserList(respAllRecords.Data.UserRecords).Render(r.Context(), w)
	case "sudo":
		respSingleRecord, err := GetSingleRecordSudo(u.userKey)
		logger.Debug(respSingleRecord)
		if err != nil {
			return err
		}
		return nil
	default:
		return server.NewUpstreamError(
			"Invalid user key", http.StatusBadRequest, nil,
		)
	}
	return nil
}
