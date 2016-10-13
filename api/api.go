package api

import (
	"github.com/gorilla/mux"
	"github.com/seenickcode/go-web-socket-chat-example/chat"
)

// API .
type API struct {
	chat *chat.Chat
}

// New initializes and wires up routes
func WireupRoutes(r *mux.Router) {
	a := &API{}
	a.chat = chat.New()

	wireupRoutes(r, a)

	return
}
