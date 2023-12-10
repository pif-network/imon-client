package task

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	. "the-gorgeouses.com/imon-client/internal/core/errors"
	"the-gorgeouses.com/imon-client/internal/core/server"
	"the-gorgeouses.com/imon-client/internal/views/components"
)

type RouterState struct {
	userKey string
	mu      sync.Mutex
}

func (rs *RouterState) SetUserKey(userKey string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.userKey = userKey
}

func (rs *RouterState) GetUserKey() string {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	return rs.userKey
}

var routerState = &RouterState{}

func PostKeyHandler(w http.ResponseWriter, r *http.Request) {
	userKey := r.PostFormValue("user-key")
	log.Printf("User key: %s", userKey)
	routerState.SetUserKey(userKey)

	resp, err := GetUserTaskLogById(userKey)
	if err != nil {
		if IsUpstreamError(err) {
			log.Printf(err.Error())
			_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
			return
		}
		log.Printf(err.Error())
		_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
		return
	}
	log.Printf("%+v\n", resp)

	res, err := GetAllUserRecords()
	if err != nil {
		if IsUpstreamError(err) {
			log.Printf(err.Error())
			_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
			return
		}
		log.Printf(err.Error())
		_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
		return
	}
	log.Printf("%+v\n", res)

	_ = CurrentTaskAndExecutionLog(resp.Data.TaskLog).Render(r.Context(), w)
	_ = ActiveUserList(res.Data.UserRecords).Render(r.Context(), w)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	state := r.PostFormValue("state")
	userKey := routerState.GetUserKey()

	err := UpdateCurrentTask(userKey, TaskState(state))
	if err != nil {
		fmt.Println(err)
		if IsUpstreamError(err) {
			log.Printf(err.Error())
			_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
			return
		}
		log.Printf(err.Error())
		_ = components.ErrorWidget(err.Error()).Render(r.Context(), w)
		return
	}
}

func GetTaskRouter() *server.Router {
	taskRouter := server.NewRouter()

	taskRouter.Post("/post-key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostKeyHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	taskRouter.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			UpdateTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return taskRouter
}
