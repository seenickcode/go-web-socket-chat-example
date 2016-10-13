package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func wireupRoutes(r *mux.Router, web *Web) {
	r.HandleFunc("/", web.index).Methods("GET")
	r.HandleFunc("/admin/users/{id}/threads/{threadID}/chat", web.chat).Methods("GET")

	fs := http.FileServer(http.Dir("./web/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	return
}
