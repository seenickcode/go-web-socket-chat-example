package api

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/seenickcode/go-web-socket-chat-example/chat"
)

func (api *API) createMessage(w http.ResponseWriter, r *http.Request) {

	userID := 0   // TODO person posting the message
	threadID := 0 // TODO thread ID of the conversation
	body := ""    // TODO text of the message

	// TODO persist your chat message here

	// broadcast to any open sockets
	broadcast := &chat.BroadcastMessage{
		DeviceIdentifier:     userID,
		ConnectionIdentifier: threadID,
		Body:                 body,
	}
	log.Debugf("sending broadcast message '%+v' to chat connection", broadcast)
	api.chat.SocketHub.Broadcast(broadcast)

	renderJSON(w, NewWrappedAPIResponse(message), http.StatusCreated)
}

func (api *API) openThreadConnection(w http.ResponseWriter, r *http.Request) {

	// get thread ID
	threadID, _ := strconv.Atoi(mux.Vars(r)["threadID"])
	if threadID <= 0 {
		log.Errorf("threadID not found in route for openChatThread")
		return
	}

	// open socket
	log.Debugf("opening chat connection for thread %v", threadID)
	ws, err := api.chat.SocketHub.SocketUpgrader.Upgrade(w, r, nil)
	if isErr(err) {
		return
	}

	// create new connection
	c := chat.NewConnection(api.chat.SocketHub, ws, threadID)
	if c == nil {
		log.Errorf("connection nil for openChatThread")
		return
	}

	// register connection with our hub, unregister when finished
	api.chat.SocketHub.RegisterCh <- c
	defer func() {
		api.chat.SocketHub.UnregisterCh <- c
	}()

	// start connection writer on another thread
	go c.WriteOutgoing()

	// start read on this thread, thus suspending this func
	c.ReadIncoming()

	log.Debugf("closing socket for %v", c)
}
