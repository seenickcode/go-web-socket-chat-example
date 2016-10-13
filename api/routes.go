package api

import "github.com/gorilla/mux"

func wireupRoutes(r *mux.Router, api *API) {
	r.HandleFunc("/api/users/{id}/threads/{threadID}/messages", api.createThreadMessage).Methods("POST")
	r.HandleFunc("/api/users/{id}/threads/{threadID}/connect", api.openThreadConnection)
	return
}
