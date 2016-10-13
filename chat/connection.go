package chat

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket     *websocket.Conn // our websocket connection
	Hub        *Hub            // parent hub
	BufferCh   chan []byte     // buffered channel of outbound messages
	Identifier int
}

func NewConnection(hub *Hub, ws *websocket.Conn, indentifier int) *Connection {
	log.Debugf("creating new connection")

	c := Connection{
		Hub:        hub,
		Socket:     ws,
		BufferCh:   make(chan []byte, 256),
		Identifier: indentifier,
	}
	return &c
}

func (c *Connection) ReadIncoming() {
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Debugf("connection reader err: %v", err)
			break
		}
		log.Debugf("connection notifying hub of message: %v", string(message))
		c.Hub.BroadcastCh <- message
	}
	log.Debugf("connection closing socket")
	c.Socket.Close()
}

func (c *Connection) WriteOutgoing() {
	// whatever comes through the buffer channel, write to our socket
	for data := range c.BufferCh {
		log.Debugf("connection writing outgoing message")

		err := c.Socket.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Debugf("connection writer err: %v", err)
			break
		}
	}
	log.Debugf("connection closing socket via writer")
	c.Socket.Close()
}
