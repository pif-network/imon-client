package task

import (
	"log"
	"net/http"

	. "the-gorgeouses.com/imon-client/internal/core/errors"
)

func PostKeyHandler(w http.ResponseWriter, r *http.Request) {
	userKey := r.PostFormValue("user-key")
	log.Printf("User key: %s", userKey)

	resp, err := GetUserTaskLogById(userKey)
	if err != nil {
		if IsUpstreamError(err) {
			log.Printf(err.Error())
			_ = ErrorWidget(err.Error()).Render(r.Context(), w)
			return
		}
		log.Printf(err.Error())
		_ = ErrorWidget(err.Error()).Render(r.Context(), w)
		return
	}
	log.Printf("%+v\n", resp)

	_ = CurrentTaskAndExecutionLog(resp.Data.TaskLog).Render(r.Context(), w)
}

func GetTaskRouter() *http.ServeMux {
	taskRouter := http.NewServeMux()

	taskRouter.HandleFunc("/post-key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostKeyHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return taskRouter
}
