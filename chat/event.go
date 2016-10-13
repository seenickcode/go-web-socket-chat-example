package chat

import (
	"fmt"
	"time"
)

type Event struct {
	Type        string    `json:"type" schema:"type"`           // type of event [followed|voted|message_sent|message_received|message_deleted]
	Time        time.Time `json:"time,omitempty" schema:"time"` // when does the event occur?
	SenderID    string    `json:"sender_id" schema:"sender_id"` // who sent the event
	SenderLat   string    `json:"sender_lat" schema:"sender_lat"`
	SenderLng   string    `json:"sender_lng" schema:"sender_lng"`
	ReceiverID  string    `json:"receiver_id" schema:"receiver_id"` // to whom
	ReceiverLat string    `json:"receiver_lat" schema:"receiver_lat"`
	ReceiverLng string    `json:"receiver_lng" schema:"receiver_lng"`
}

func (e *Event) Validate() error {
	if e.Type != "followed" &&
		e.Type != "voted" &&
		e.Type != "message_sent" &&
		e.Type != "message_received" &&
		e.Type != "message_deleted" {
		return fmt.Errorf("Invalid event type '%v'", e.Type)
	}
	return nil
}

// sort.Interface

type EventsByTime []*Event

func (s EventsByTime) Len() int {
	return len(s)
}

func (s EventsByTime) Less(i, j int) bool {
	return s[i].Time.Before(s[j].Time)
}

func (s EventsByTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
