package hbslog

import (
	"log/slog"
	"testing"

	"github.com/honeybadger-io/honeybadger-go"
)

func newTestClient() (*honeybadger.Client, *honeybadger.TestBackend) {
	backend := &honeybadger.TestBackend{}
	client := honeybadger.New(honeybadger.Configuration{
		APIKey:  "test-key",
		Backend: backend,
		Sync:    true,
	})
	return client, backend
}

func TestWithAttrs(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test")

	baseAttrs := []slog.Attr{
		slog.String("service", "api"),
		slog.Int("version", 1),
	}
	loggerWithAttrs := slog.New(handler.WithAttrs(baseAttrs))

	loggerWithAttrs.Info("test message", "request_id", "123")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "test" {
		t.Errorf("expected event type 'test', got %q", event.EventType)
	}

	if event.Data["level"] != "INFO" {
		t.Errorf("expected level='INFO', got %v", event.Data["level"])
	}
	if event.Data["message"] != "test message" {
		t.Errorf("expected message='test message', got %v", event.Data["message"])
	}
	if event.Data["service"] != "api" {
		t.Errorf("expected service='api', got %v", event.Data["service"])
	}
	if event.Data["version"] != int64(1) {
		t.Errorf("expected version=1, got %v (type %T)", event.Data["version"], event.Data["version"])
	}
	if event.Data["request_id"] != "123" {
		t.Errorf("expected request_id='123', got %v", event.Data["request_id"])
	}
}

func TestWithAttrsChaining(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test")

	handler1 := handler.WithAttrs([]slog.Attr{slog.String("service", "api")})
	handler2 := handler1.WithAttrs([]slog.Attr{slog.Int("version", 2)})

	logger := slog.New(handler2)
	logger.Info("chained message")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Data["service"] != "api" {
		t.Errorf("expected service='api', got %v", event.Data["service"])
	}
	if event.Data["version"] != int64(2) {
		t.Errorf("expected version=2, got %v", event.Data["version"])
	}
}

func TestWithGroup(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test")

	loggerWithGroup := slog.New(handler.WithGroup("request"))
	loggerWithGroup.Info("test message", "id", "123", "method", "GET")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Data["level"] != "INFO" {
		t.Errorf("expected level='INFO', got %v", event.Data["level"])
	}
	if event.Data["message"] != "test message" {
		t.Errorf("expected message='test message', got %v", event.Data["message"])
	}

	requestGroup, ok := event.Data["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'request' to be a map, got %T", event.Data["request"])
	}

	if requestGroup["id"] != "123" {
		t.Errorf("expected request.id='123', got %v", requestGroup["id"])
	}
	if requestGroup["method"] != "GET" {
		t.Errorf("expected request.method='GET', got %v", requestGroup["method"])
	}
}

func TestWithGroupNested(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test")

	handler1 := handler.WithGroup("request")
	handler2 := handler1.WithGroup("headers")

	logger := slog.New(handler2)
	logger.Info("nested group", "content-type", "application/json")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	requestGroup, ok := event.Data["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'request' to be a map, got %T", event.Data["request"])
	}

	headersGroup, ok := requestGroup["headers"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'request.headers' to be a map, got %T", requestGroup["headers"])
	}

	if headersGroup["content-type"] != "application/json" {
		t.Errorf("expected content-type='application/json', got %v", headersGroup["content-type"])
	}
}

func TestWithAttrsAndGroup(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test")

	handler1 := handler.WithAttrs([]slog.Attr{slog.String("service", "api")})
	handler2 := handler1.WithGroup("request")

	logger := slog.New(handler2)
	logger.Info("mixed", "id", "456")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Data["service"] != "api" {
		t.Errorf("expected service='api', got %v", event.Data["service"])
	}

	requestGroup, ok := event.Data["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'request' to be a map, got %T", event.Data["request"])
	}

	if requestGroup["id"] != "456" {
		t.Errorf("expected request.id='456', got %v", requestGroup["id"])
	}
}

func TestWithAttrsBeforeAndAfterGroup(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	h1 := h.WithAttrs([]slog.Attr{slog.String("service", "api")})
	h2 := h1.WithGroup("http").WithAttrs([]slog.Attr{slog.String("method", "GET")})

	slog.New(h2).Info("x")
	e := backend.GetEvents()[0]

	if got := e.Data["service"]; got != "api" {
		t.Fatalf("top-level pre-attr lost: %v", got)
	}
	httpM := e.Data["http"].(map[string]any)
	if httpM["method"] != "GET" {
		t.Fatalf("grouped pre-attr missing: %v", httpM["method"])
	}
}

func TestWithEventType(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client)

	handler2 := handler.WithEventType("audit")
	logger := slog.New(handler2)
	logger.Info("audit event", "user_id", "123")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "audit" {
		t.Errorf("expected event type 'audit', got %q", event.EventType)
	}
	if event.Data["user_id"] != "123" {
		t.Errorf("expected user_id='123', got %v", event.Data["user_id"])
	}
}

func TestWithEventTypePreservesAttrsAndGroups(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client)

	handler2 := handler.WithEventType("api_call").
		WithAttrs([]slog.Attr{slog.String("env", "prod")}).
		WithGroup("request")

	logger := slog.New(handler2)
	logger.Info("api call", "path", "/users")

	events := backend.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventType != "api_call" {
		t.Errorf("expected event type 'api_call', got %q", event.EventType)
	}
	if event.Data["env"] != "prod" {
		t.Errorf("expected env='prod', got %v", event.Data["env"])
	}

	requestGroup, ok := event.Data["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'request' to be a map, got %T", event.Data["request"])
	}
	if requestGroup["path"] != "/users" {
		t.Errorf("expected request.path='/users', got %v", requestGroup["path"])
	}
}

func TestInlineSlogGroup(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")
	slog.New(h.WithGroup("http")).Info("x",
		slog.Group("resp", slog.Int("status", 200), slog.String("ct", "json")),
	)
	e := backend.GetEvents()[0]
	http := e.Data["http"].(map[string]any)
	resp := http["resp"].(map[string]any)
	if resp["status"] != int64(200) || resp["ct"] != "json" {
		t.Fatalf("inline group not expanded: %#v", resp)
	}
}

type counterValuer struct {
	v *int
}

func (c counterValuer) LogValue() slog.Value {
	return slog.IntValue(*c.v)
}

func TestLogValuerResolvesAtHandleTime(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	v := 0
	lv := slog.Any("counter", counterValuer{&v})
	l := slog.New(h.WithAttrs([]slog.Attr{lv}))

	v = 1
	l.Info("first")
	v = 2
	l.Info("second")

	ev := backend.GetEvents()
	if ev[0].Data["counter"] != int64(1) || ev[1].Data["counter"] != int64(2) {
		t.Fatalf("LogValuer not resolved at handle time: %+v", []any{ev[0].Data["counter"], ev[1].Data["counter"]})
	}
}

func TestEventTypeCanBeOverridden(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	slog.New(h).Info("msg", "event_type", "user_value")

	e := backend.GetEvents()[0]
	if e.EventType != "user_value" {
		t.Errorf("event_type should be overridable per-call, got %v", e.EventType)
	}
}

func TestUserCanSetTimestamp(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	slog.New(h).Info("msg", "ts", "2024-01-01T00:00:00Z")

	e := backend.GetEvents()[0]
	if e.Data["ts"] != "2024-01-01T00:00:00Z" {
		t.Errorf("user should be able to set ts, got %v", e.Data["ts"])
	}
}

// 1) Immutability: WithGroup/WithAttrs must not mutate the parent handler/logger.
func TestImmutability(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	parent := slog.New(h)
	child := slog.New(h.WithGroup("http").WithAttrs([]slog.Attr{slog.String("a", "b")}))

	parent.Info("base")
	child.Info("child", "k", "v")

	evs := backend.GetEvents()
	if len(evs) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}
	if _, ok := evs[0].Data["http"]; ok {
		t.Fatal("parent mutated by WithGroup/WithAttrs")
	}
	if _, ok := evs[1].Data["http"]; !ok {
		t.Fatal("child missing group namespace")
	}
}

//  2. Depth semantics: pre-attrs added before a group stay top-level;
//     pre-attrs added after a group live under that group.
func TestPreAttrsDepthBeforeAfterGroup(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("test")

	h1 := h.WithAttrs([]slog.Attr{slog.String("service", "api")}) // depth=0 → top-level
	h2 := h1.WithGroup("http").WithAttrs([]slog.Attr{slog.String("method", "GET")})

	slog.New(h2).Info("msg")

	ev := backend.GetEvents()[0]
	if got := ev.Data["service"]; got != "api" {
		t.Fatalf("top-level pre-attr lost: %v", got)
	}
	httpM, ok := ev.Data["http"].(map[string]any)
	if !ok {
		t.Fatalf("'http' should be a map, got %T", ev.Data["http"])
	}
	if httpM["method"] != "GET" {
		t.Fatalf("grouped pre-attr missing: %v", httpM["method"])
	}
}

//  3. Envelope control: event_type can come from handler default (WithEventType)
//     but can be overridden per-call via the event_type attribute.
func TestEventTypeDefaultWithOverride(t *testing.T) {
	client, backend := newTestClient()
	h := New(client).WithEventType("audit") // set default

	logger := slog.New(h)
	logger.Info("uses default")                            // uses "audit"
	logger.Info("overrides", "event_type", "custom_type") // overrides to "custom_type"

	evs := backend.GetEvents()
	if len(evs) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evs))
	}

	if evs[0].EventType != "audit" {
		t.Errorf("expected first event type 'audit', got %q", evs[0].EventType)
	}
	if evs[1].EventType != "custom_type" {
		t.Errorf("expected second event type 'custom_type', got %q", evs[1].EventType)
	}
}

func TestLogLevelFiltering(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client).WithEventType("test").WithLevel(slog.LevelWarn)

	logger := slog.New(handler)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	events := backend.GetEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 events (warn and error), got %d", len(events))
	}

	if events[0].Data["level"] != "WARN" {
		t.Errorf("expected first event level='WARN', got %v", events[0].Data["level"])
	}
	if events[1].Data["level"] != "ERROR" {
		t.Errorf("expected second event level='ERROR', got %v", events[1].Data["level"])
	}
}

func TestLogLevelVar(t *testing.T) {
	client, backend := newTestClient()
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelWarn)

	handler := New(client).WithEventType("test").WithLevel(levelVar)
	logger := slog.New(handler)

	logger.Info("info message")
	logger.Warn("warn message")

	if len(backend.GetEvents()) != 1 {
		t.Fatalf("expected 1 event before level change, got %d", len(backend.GetEvents()))
	}

	levelVar.Set(slog.LevelDebug)
	logger.Info("second info message")

	events := backend.GetEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 events after level change, got %d", len(events))
	}
}

func TestEventTypePerLogCall(t *testing.T) {
	client, backend := newTestClient()
	handler := New(client)
	logger := slog.New(handler)

	logger.Info("user signup", "event_type", "user_lifecycle", "user_id", 123)
	logger.Info("payment processed", "event_type", "payment", "amount", 99.99)
	logger.Info("regular log message", "foo", "bar")

	events := backend.GetEvents()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	if events[0].EventType != "user_lifecycle" {
		t.Errorf("expected first event type 'user_lifecycle', got %q", events[0].EventType)
	}
	if events[1].EventType != "payment" {
		t.Errorf("expected second event type 'payment', got %q", events[1].EventType)
	}
	if events[2].EventType != "log" {
		t.Errorf("expected third event type 'log' (default), got %q", events[2].EventType)
	}

	// event_type in data should match the EventType envelope
	if events[0].Data["event_type"] != "user_lifecycle" {
		t.Errorf("expected data event_type 'user_lifecycle', got %v", events[0].Data["event_type"])
	}
	if events[1].Data["event_type"] != "payment" {
		t.Errorf("expected data event_type 'payment', got %v", events[1].Data["event_type"])
	}
}
