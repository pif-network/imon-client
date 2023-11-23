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
			log.Printf("Failed to get user task log.")
			log.Printf(err.Error())
			_ = ErrorWidget(err.Error()).Render(r.Context(), w)
			return
		}
		log.Printf("ailed to get user task log.")
		log.Printf(err.Error())
		_ = ErrorWidget(err.Error()).Render(r.Context(), w)
		return
	}
	log.Printf("%+v\n", resp)

	// if err := templates.ExecuteTemplate(w, "task-list", resp); err != nil {
	// 	log.Printf("Failed to execute template.")
	// 	log.Printf(err.Error())
	// }

	_ = CurrentTaskAndExecutionLog(resp.Data.TaskLog).Render(r.Context(), w)
}
