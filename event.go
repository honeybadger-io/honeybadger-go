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

func newEventPayload(eventData map[string]interface{}) *eventPayload {
	data := make(map[string]interface{})
	for k, v := range eventData {
		data[k] = v
	}
	
	data["ts"] = time.Now().UTC().Format(time.RFC3339)
	
	return &eventPayload{data: data}
}