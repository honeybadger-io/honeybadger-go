package honeybadger

import (
	"maps"
	"time"
)

type eventPayload struct {
	data map[string]any
}

func (e *eventPayload) toJSON() []byte {
	h := hash(e.data)
	return h.toJSON()
}

func newEventPayload(eventType string, eventData map[string]any) *eventPayload {
	data := make(map[string]any)
	maps.Copy(data, eventData)

	data["event_type"] = eventType
	if _, ok := data["ts"]; !ok {
		data["ts"] = time.Now().UTC().Format(time.RFC3339)
	}

	return &eventPayload{data: data}
}

type eventBatch struct {
	events []*eventPayload
}

func (b *eventBatch) toJSON() []byte {
	var events []map[string]any
	for _, event := range b.events {
		events = append(events, event.data)
	}
	h := hash{"events": events}
	return h.toJSON()
}

func newEventBatch(events []*eventPayload) *eventBatch {
	return &eventBatch{events: events}
}
