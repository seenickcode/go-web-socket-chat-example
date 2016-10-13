package chat

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	log "github.com/Sirupsen/logrus"
)

type Hub struct {
	SocketUpgrader websocket.Upgrader // allows us to generate new web socket connections
	Connections    map[*Connection]bool
	BroadcastCh    chan []byte      // we send data here
	RegisterCh     chan *Connection // register requests from the connections
	UnregisterCh   chan *Connection // unregister requests from connections
}

func NewHub() *Hub {
	h := &Hub{
		SocketUpgrader: websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		Connections:    make(map[*Connection]bool),
		BroadcastCh:    make(chan []byte),
		RegisterCh:     make(chan *Connection),
		UnregisterCh:   make(chan *Connection),
	}
	return h
}

type BroadcastMessage struct {
	ConnectionIdentifier int    `json:"connection_identifier"`
	DeviceIdentifier     int    `json:"device_identifier"`
	Body                 string `json:"body"`
}

func (h *Hub) Run() {
	log.Infof("hub starting, please wait...")
	for {
		select {
		case c := <-h.RegisterCh:

			log.Infof("hub registering new connection")

			h.Connections[c] = true

			log.Infof("hub has %d current connections", len(h.Connections))

		case c := <-h.UnregisterCh:

			log.Infof("hub unregistering connection %v", c)

			delete(h.Connections, c)
			close(c.BufferCh)
			log.Debugf("hub now has %d connections", len(h.Connections))
		}
	}
}

// Broadcast broadcasts message to connections for specific identifier
func (h *Hub) Broadcast(message *BroadcastMessage) {
	log.Infof("hub broadcasting %+v", message)

	data, err := json.Marshal(message)
	if err != nil {
		log.Errorf("couldn't marshal message %+v: %v", message, err)
		return
	}
	log.Debugf("marshalled message for broadcasting %+v", message)

	for c := range h.Connections {
		// send message to all relevant connections with same identifier
		if c.Identifier == message.ConnectionIdentifier {
			c.BufferCh <- data
		}
	}
}
