package infrastruction

import (
	"log"
	"time"
)

type Event struct {
	RelationID  int64
	ReferenceID int64
	Tags        string
	Detail      string
	Time        time.Time
}

type EventListener struct {
	EventChannel chan Event
}

var EventListenerInstance *EventListener

func init() {
	EventListenerInstance = &EventListener{
		EventChannel: make(chan Event),
	}

	go func() {
		for event := range EventListenerInstance.EventChannel {
			log.Println("event received: ", event)
		}
	}()
}
