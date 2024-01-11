package task

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	user User
	mu   sync.Mutex
}

func (rs *RouterState) SetUserKey(userKey string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// $type:$name:$id
	parts := strings.Split(userKey, ":")
	if len(parts) != 3 {
		// Invalid, but not handling here.
		rs.user = User{
			userKey: userKey,
		}
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		// Invalid, but not handling here.
		rs.user = User{
			userKey: userKey,
		}
		return
	}

	rs.user = User{
		userKey:  userKey,
		userType: parts[0],
		name:     parts[1],
		id:       id,
	}
}

func (rs *RouterState) GetUserKey() User {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	return rs.user
}

var routerState = &RouterState{}

func PostKeyHandler(w http.ResponseWriter, r *http.Request) {
	userKey := r.PostFormValue("user-key")
	logger.Info("", "userKey", userKey)
	routerState.SetUserKey(userKey)

	w.Header().Set("HX-Trigger", shared.ClientEvt.ShouldRefresh)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	state := r.PostFormValue("state")
	user := routerState.GetUserKey()

	err := UpdateCurrentTask(user.userKey, TaskState(state))
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

	w.Header().Set("HX-Trigger", shared.ClientEvt.ShouldRefresh)
}

func RefreshAppDataHandler(w http.ResponseWriter, r *http.Request) {
	user := routerState.GetUserKey()

	err := user.RefreshData(w, r)
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
