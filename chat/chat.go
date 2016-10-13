package chat

import _ "github.com/lib/pq" // postgres

// Chat provides the ability to for clients to send and receive chat events
type Chat struct {
	SocketHub *Hub
}

// New creates a new Chat instance
func New() *Chat {
	c := &Chat{}
	c.SocketHub = NewHub()
	go c.SocketHub.Run() // start listening/broadcasting socket messages
	return c
}

// Close closes all connections
func (c *Chat) Close() {
}
