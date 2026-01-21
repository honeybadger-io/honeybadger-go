package hbzerolog

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/honeybadger-io/honeybadger-go"
	"github.com/rs/zerolog"
)

type Writer struct {
	c         *honeybadger.Client
	eventType string
	timeKey   string
	levelKey  string
}

type Option func(*Writer)

func WithEventType(t string) Option { return func(w *Writer) { w.eventType = t } }
func WithKeys(timeKey, levelKey string) Option {
	return func(w *Writer) { w.timeKey, w.levelKey = timeKey, levelKey }
}

func New(c *honeybadger.Client, opts ...Option) *Writer {
	w := &Writer{
		c:         c,
		eventType: "log",
		timeKey:   "time",
		levelKey:  "level",
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

func (w *Writer) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(p))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		m = map[string]any{"message": string(bytes.TrimSpace(p))}
	}

	eventType := w.eventType
	if et, ok := m["event_type"].(string); ok && et != "" {
		eventType = et
		delete(m, "event_type")
	}

	if _, ok := m[w.levelKey]; !ok {
		if level == zerolog.NoLevel {
			level = zerolog.InfoLevel
		}
		m[w.levelKey] = level.String()
	}
	if _, ok := m[w.timeKey]; !ok {
		m[w.timeKey] = time.Now().UTC().Format(time.RFC3339Nano)
	}

	if w.timeKey != "ts" {
		if ts, ok := m[w.timeKey]; ok {
			m["ts"] = ts
			delete(m, w.timeKey)
		}
	}

	if err := w.c.Event(eventType, m); err != nil {
		if w.c != nil && w.c.Config != nil && w.c.Config.Logger != nil {
			w.c.Config.Logger.Printf("zerolog writer failed to send event: %v\n", err)
		}
		return len(p), err
	}
	return len(p), nil
}

func (w *Writer) Write(p []byte) (int, error) { return w.WriteLevel(zerolog.NoLevel, p) }
