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
