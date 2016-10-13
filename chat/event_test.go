package chat

import (
	"io/ioutil"
	"net/http"
	"sort"
	"testing"
	"time"
)

var validateableEvents map[*Event]bool
var unsortedEvents []*Event

func TestValidateEvent(t *testing.T) {
	setup()
	defer teardown()

	for e, expectedValid := range validateableEvents {
		t.Logf("validating event: %+v", e)
		err := e.Validate()
		if expectedValid && err != nil {
			t.Errorf("expected valid value of %v, got %v, error was: %v", expectedValid, !expectedValid, err)
		}
	}
}

func TestSortableEvents(t *testing.T) {
	setup()
	defer teardown()

	// sort events
	events := []*Event{}
	copy(unsortedEvents, events)
	sort.Sort(EventsByTime(events))

	// ensure sorted
	for ndx, e := range events {
		if ndx != len(events)-1 {
			nextTime := events[ndx+1].Time
			if e.Time.After(nextTime) {
				t.Errorf("event time is not sorted, got %v and expected to be before %v", e.Time, nextTime)
			}
		}
	}
}

func BenchmarkValidateEvent(b *testing.B) {

	for n := 0; n < b.N; n++ {
		for e, _ := range validateableEvents {
			e.Validate()
		}
	}
}

func setup() {
	// map of events and "is valid"
	validateableEvents = map[*Event]bool{
		&Event{
			Type:        "message_sent",
			SenderID:    "217152a577b7b013b4636e203d4dea9a",
			SenderLat:   "44.34",
			SenderLng:   "17.35",
			ReceiverID:  "646152a577b7b013b4636e203d4dea9b",
			ReceiverLat: "44.1011",
			ReceiverLng: "17.9311",
		}: true,
		&Event{
			Type: "invalid_type_123",
		}: false,
	}

	// unsorted events
	now := time.Now()
	unsortedEvents = []*Event{
		&Event{
			Type: "message_sent",
			Time: now.Add(time.Second * 3),
		},
		&Event{
			Type: "message_sent",
			Time: now.Add(time.Second * 1),
		},
		&Event{
			Type: "message_sent",
			Time: now.Add(time.Second * 2),
		},
		&Event{
			Type: "message_sent",
			Time: now.Add(time.Second * 7),
		},
		&Event{
			Type: "message_sent",
			Time: now.Add(time.Second * 5),
		},
	}

}

func teardown() {
}

func httpResponseToString(r *http.Response) string {
	data, _ := ioutil.ReadAll(r.Body)
	return string(data)
}
