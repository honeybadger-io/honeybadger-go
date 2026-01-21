package honeybadger

import "sync"

type TestBackend struct {
	Events []EventData
	mu     sync.Mutex
}

type EventData struct {
	EventType string
	Data      map[string]any
}

func (b *TestBackend) Notify(_ Feature, _ Payload) error {
	return nil
}

func (b *TestBackend) Event(events []*eventPayload) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, e := range events {
		eventType, _ := e.data["event_type"].(string)
		b.Events = append(b.Events, EventData{
			EventType: eventType,
			Data:      e.data,
		})
	}
	return nil
}

func (b *TestBackend) GetEvents() []EventData {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Events
}
