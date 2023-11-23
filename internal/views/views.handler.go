package views

import (
	"net/http"

	"the-gorgeouses.com/imon-client/internal/views/pages"
)

func ServeRootView(w http.ResponseWriter, r *http.Request) {
	_ = pages.Index().Render(r.Context(), w)
}
