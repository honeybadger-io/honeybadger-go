package honeybadger

import (
	"time"
)

type eventPayload struct {
	data map[string]interface{}
}

func (e *eventPayload) toJSON() []byte {
	h := hash(e.data)
	return h.toJSON()
}

func newEventPayload(eventType string, eventData map[string]interface{}) *eventPayload {
	data := make(map[string]interface{})
	for k, v := range eventData {
		data[k] = v
	}
	
	data["event_type"] = eventType
	data["ts"] = time.Now().UTC().Format(time.RFC3339)
	
	return &eventPayload{data: data}
}

type eventBatch struct {
	events []*eventPayload
}

func (b *eventBatch) toJSON() []byte {
	var events []map[string]interface{}
	for _, event := range b.events {
		events = append(events, event.data)
	}
	h := hash{"events": events}
	return h.toJSON()
}

func newEventBatch(events []*eventPayload) *eventBatch {
	return &eventBatch{events: events}
}