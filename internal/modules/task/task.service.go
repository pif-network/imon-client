package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/core/server"
)

type UserTaskLogResponse struct {
	Data struct {
		TaskLog Record `json:"task_log"`
	} `json:"data"`
	Status string `json:"status"`
}

func GetUserRecord(userKey string) (UserTaskLogResponse, error) {
	payload := generatePayload(userKey, "user", UpstreamEventType.GetSingleRecord)
	logger.Debug(payload)
	resp, err := http.Post(
		"http://localhost:8000/v1/rpc/user",
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
		UserRecords []Record `json:"user_records"`
	} `json:"data"`
	Status string `json:"status"`
}

func GetAllUserRecords(userKey string) (AllUserRecordsResponse, error) {
	payload := generatePayload(userKey, "user", UpstreamEventType.GetAllUserRecords)
	logger.Debug(payload)
	resp, err := http.Post(
		"http://localhost:8000/v1/rpc/user",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
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

func generatePayload(userKey string, userType string, eventType string) string {
	payload := RpcPayload{
		Metadata: struct {
			Of string `json:"of"`
		}{
			Of: userType,
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

func GetSingleRecordSudo(userKey string, userType string) (SingleRecordResponse, error) {
	payload := generatePayload(userKey, userType, UpstreamEventType.GetSingleRecord)
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
		var wg sync.WaitGroup
		wg.Add(2)

		errCtx, cancel := context.WithCancel(r.Context())
		defer cancel()

		respSingleRecordCh := make(chan UserTaskLogResponse, 1)
		respAllRecordsCh := make(chan AllUserRecordsResponse, 1)
		errCh := make(chan error, 1)

		go func() {
			defer wg.Done()
			select {
			case <-errCtx.Done():
				return
			default:
			}
			respTaskLog, err := GetUserRecord(u.userKey)
			if err != nil {
				cancel()
				close(respSingleRecordCh)
				logger.Debug(err.Error())
				errCh <- err
				return
			}
			respSingleRecordCh <- respTaskLog
			logger.Debug("Got single record", "respTaskLog", respTaskLog)
		}()
		go func() {
			defer wg.Done()
			select {
			case <-errCtx.Done():
				return
			default:
			}
			respAllRecords, err := GetAllUserRecords(u.userKey)
			logger.Debug("Got all records", "respAllRecords", respAllRecords)
			if err != nil {
				cancel()
				close(respAllRecordsCh)
				logger.Debug(err.Error())
				errCh <- err
				return
			}
			respAllRecordsCh <- respAllRecords
		}()

		wg.Wait()

		if errCtx.Err() != nil {
			return <-errCh
		}
		close(errCh)

		respTaskLog := <-respSingleRecordCh
		logger.Debug("Got single record", "respTaskLog", respTaskLog)
		_ = CurrentTaskAndExecutionLog(respTaskLog.Data.TaskLog).Render(r.Context(), w)
		respAllRecords := <-respAllRecordsCh
		logger.Debug("Got all records", "respAllRecords", respAllRecords)
		_ = ActiveUserList(respAllRecords.Data.UserRecords).Render(r.Context(), w)

		// select {
		// case err := <-errCh:
		// 	return err
		// default:
		// 	respTaskLog := <-respSingleRecordCh
		// 	logger.Debug("Got single record", "respTaskLog", respTaskLog)
		// 	_ = CurrentTaskAndExecutionLog(respTaskLog.Data.TaskLog).Render(r.Context(), w)
		// 	respAllRecords := <-respAllRecordsCh
		// 	logger.Debug("Got all records", "respAllRecords", respAllRecords)
		// 	_ = ActiveUserList(respAllRecords.Data.UserRecords).Render(r.Context(), w)
		// }
	case "sudo":
		respSingleRecord, err := GetSingleRecordSudo(u.userKey, u.userType)
		if err != nil {
			return err
		}
		logger.Debug(respSingleRecord)
		_ = STaskList(respSingleRecord.Data.PublishedTasks).Render(r.Context(), w)
		return nil
	default:
		return server.NewUpstreamError(
			"Invalid user key", http.StatusBadRequest, nil,
		)
	}
	return nil
}
