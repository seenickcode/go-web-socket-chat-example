package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (web *Web) chat(w http.ResponseWriter, r *http.Request) {
	threadID, _ := strconv.Atoi(mux.Vars(r)["threadID"])
	data := map[string]interface{}{
		"ThreadID": threadID,
	}
	renderPageWithLayout(web, w, r, "admin", "chat.html", "admin.html", data)
}
