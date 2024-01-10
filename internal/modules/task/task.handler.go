package task

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"

	"the-gorgeouses.com/imon-client/internal/core/server"
	"the-gorgeouses.com/imon-client/internal/core/shared"
	"the-gorgeouses.com/imon-client/internal/views/components"
)

var CompSwapId = struct {
	KeyForm string
}{
	KeyForm: "key-form-error",
}

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
	log.Info("", "userKey", userKey)
	routerState.SetUserKey(userKey)

	w.Header().Set("HX-Trigger", shared.ClientEvt.ShouldRefresh)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	state := r.PostFormValue("state")
	userKey := routerState.GetUserKey()

	err := UpdateCurrentTask(userKey, TaskState(state))
	if err != nil {
		logger.Error(err.Error())

		if ferr, ok := server.FixableByClient(err); ok {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				ferr.Display()).Render(r.Context(),
				w,
			)
			return
		} else {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				"Internal_Error: Well, something is broken.").Render(r.Context(),
				w,
			)
			return
		}
	}

	w.Header().Set("HX-Trigger", "task_updated")
}

func RefreshAppDataHandler(w http.ResponseWriter, r *http.Request) {
	userKey := routerState.GetUserKey()
	if userKey == "" {
		http.Error(w, "Incorrect flow - First post the user key.", http.StatusForbidden)
		return
	}

	respTaskLog, err := GetUserTaskLogById(userKey)
	if err != nil {
		logger.Error(err.Error())

		if ferr, ok := server.FixableByClient(err); ok {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				ferr.Display()).Render(r.Context(),
				w,
			)
			return
		} else {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				"Internal_Error: Well, something is broken.").Render(r.Context(),
				w,
			)
			return
		}
	}

	respAllRecords, err := GetAllUserRecords()
	if err != nil {
		logger.Error(err.Error())

		if ferr, ok := server.FixableByClient(err); ok {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				ferr.Display()).Render(r.Context(),
				w,
			)
			return
		} else {
			_ = components.ErrorWidget(
				CompSwapId.KeyForm,
				"Internal_Error: Well, something is broken.").Render(r.Context(),
				w,
			)
			return
		}
	}

	_ = CurrentTaskAndExecutionLog(respTaskLog.Data.TaskLog).Render(r.Context(), w)
	_ = ActiveUserList(respAllRecords.Data.UserRecords).Render(r.Context(), w)
}

func GetTaskRouter() *server.Router {
	taskRouter := server.NewRouter()

	taskRouter.Post("/post-key", func(w http.ResponseWriter, r *http.Request) {
		PostKeyHandler(w, r)
	})
	taskRouter.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		UpdateTaskHandler(w, r)
	})
	taskRouter.Get("/refresh", func(w http.ResponseWriter, r *http.Request) {
		RefreshAppDataHandler(w, r)
	})

	return taskRouter
}
