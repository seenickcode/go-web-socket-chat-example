package api

import (
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent"
)

func wireupRoutes(r *mux.Router, api *API) {
	r.HandleFunc(newrelic.WrapHandleFunc(*api.nr, "/api/users/{id}/threads/{threadID}/messages", api.createMessage)).Methods("POST")
	r.HandleFunc(newrelic.WrapHandleFunc(*api.nr, "/api/users/{id}/threads/{threadID}/connect", api.openThreadConnection))
	return
}
