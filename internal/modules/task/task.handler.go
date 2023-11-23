package task

import (
	"log"
	"net/http"

	. "the-gorgeouses.com/imon-client/internal/core"
	"the-gorgeouses.com/imon-client/internal/views/pages"
)

func ServeRootView(w http.ResponseWriter, r *http.Request) {
	_ = pages.Index().Render(r.Context(), w)
}

func PostKeyHandler(w http.ResponseWriter, r *http.Request) {
	userKey := r.PostFormValue("user-key")
	log.Printf("User key: %s", userKey)

	resp, err := GetUserTaskLogById(userKey)
	if err != nil {
		if IsUpstreamError(err) {
			log.Printf("Failed to get user task log.")
			log.Printf(err.Error())
			// templates.ExecuteTemplate(w, "invalid-user-key", "Cannot reach service.")
			return
		}
		log.Printf("ailed to get user task log.")
		log.Printf(err.Error())
		// templates.ExecuteTemplate(w, "invalid-user-key", "Cannot reach service.")
		return
	}
	log.Printf("%+v\n", resp)

	// if err := templates.ExecuteTemplate(w, "task-list", resp); err != nil {
	// 	log.Printf("Failed to execute template.")
	// 	log.Printf(err.Error())
	// }
}
