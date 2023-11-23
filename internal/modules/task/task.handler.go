package task

import (
	"log"
	"net/http"

	. "the-gorgeouses.com/imon-client/internal/core"
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
