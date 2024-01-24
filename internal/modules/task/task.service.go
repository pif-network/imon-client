package task

import (
	"bytes"
	"context"
	"encoding/json"
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
	payload := generatePayload(userKey, "user", UpstreamEventType.GetSingleRecord, nil)
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
	payload := generatePayload(userKey, "user", UpstreamEventType.GetAllUserRecords, nil)
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
	payload := generatePayload(userKey, "user", UpstreamEventType.UpdateTask, map[string]interface{}{
		"state": taskState,
	})
	logger.Debug(payload)

	resp, err := http.Post(
		"http://localhost:8000/v1/rpc/user",
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

func generatePayload(userKey string, userType string, eventType string, additionalPayload map[string]interface{}) string {
	var payload RpcPayload
	if userKey == "" {
		payload = RpcPayload{
			Metadata: struct {
				Of string `json:"of"`
			}{
				Of: userType,
			},
			Payload: map[string]interface{}{
				"event_type": eventType,
			},
		}

	} else {
		payload = RpcPayload{
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
	}
	if additionalPayload != nil {
		for k, v := range additionalPayload {
			payload.Payload[k] = v
		}
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
	payload := generatePayload(userKey, userType, UpstreamEventType.GetSingleRecord, nil)
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
		// NOTE: The 2 goroutines fire at the same time, so no need to check for
		// context cancellation at the beginning of each goroutine.
		// If goroutine 1 finishes first and causes an error, `errCtx` will be
		// cancelled, and no matter the status of goroutine 2, it will just return.

		var wg sync.WaitGroup
		wg.Add(2)

		errCtx, cancel := context.WithCancel(r.Context())
		defer cancel()

		respSingleRecordCh := make(chan UserTaskLogResponse, 1)
		respAllRecordsCh := make(chan AllUserRecordsResponse, 1)
		errCh := make(chan error, 1)

		go func() {
			defer wg.Done()
			respTaskLog, err := GetUserRecord(u.userKey)
			if err != nil {
				select {
				case <-errCtx.Done():
					logger.Debug("Had error")
					return
				default:
					cancel()
					logger.Debug(err.Error())
					errCh <- err
					return
				}
			}
			select {
			case <-errCtx.Done():
				return
			case respSingleRecordCh <- respTaskLog:
				logger.Debug("Got single record", "respTaskLog", respTaskLog)
			default:
			}
		}()
		go func() {
			defer wg.Done()
			respAllRecords, err := GetAllUserRecords(u.userKey)
			if err != nil {
				select {
				case <-errCtx.Done():
					logger.Debug("Had error")
					return
				default:
					cancel()
					logger.Debug(err.Error())
					errCh <- err
					return
				}
			}
			select {
			case <-errCtx.Done():
				return
			case respAllRecordsCh <- respAllRecords:
				logger.Debug("Got all records", "respAllRecords", respAllRecords)
			default:
			}
		}()

		wg.Wait()
		close(respSingleRecordCh)
		close(respAllRecordsCh)

		if errCtx.Err() != nil {
			return <-errCh
		}
		close(errCh)

		respTaskLog := <-respSingleRecordCh
		logger.Debug("RENDERING single record")
		_ = CurrentTaskAndExecutionLog(respTaskLog.Data.TaskLog).Render(r.Context(), w)
		respAllRecords := <-respAllRecordsCh
		logger.Debug("RENDERING all records")
		_ = ActiveUserList(respAllRecords.Data.UserRecords).Render(r.Context(), w)
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
