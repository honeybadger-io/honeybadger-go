package honeybadger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/mock"
)

type eventReqResp struct {
	Body     []byte
	Response chan int
}

var (
	mux           *http.ServeMux
	ts            *httptest.Server
	requests      []*HTTPRequest
	defaultConfig = *Config
)

type MockedHandler struct {
	mock.Mock
}

func (h *MockedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Called()
}

type HTTPRequest struct {
	Request *http.Request
	Body    []byte
}

func (h *HTTPRequest) decodeJSON() hash {
	var dat hash
	err := json.Unmarshal(h.Body, &dat)
	if err != nil {
		panic(err)
	}
	return dat
}

func newHTTPRequest(r *http.Request) *HTTPRequest {
	body, _ := io.ReadAll(r.Body)
	return &HTTPRequest{r, body}
}

func assertMethod(t *testing.T, r *http.Request, method string) {
	if r.Method != method {
		t.Errorf("Unexpected request method. actual=%#v expected=%#v", r.Method, method)
	}
}

func setup(t *testing.T) {
	mux = http.NewServeMux()
	ts = httptest.NewServer(mux)
	requests = []*HTTPRequest{}
	mux.HandleFunc("/v1/notices",
		func(w http.ResponseWriter, r *http.Request) {
			assertMethod(t, r, "POST")
			requests = append(requests, newHTTPRequest(r))
			w.WriteHeader(201)
			fmt.Fprint(w, `{"id":"87ded4b4-63cc-480a-b50c-8abe1376d972"}`)
		},
	)

	*DefaultClient.Config = *newConfig(Configuration{APIKey: "badgers", Endpoint: ts.URL})
}

func setupEvents(t *testing.T) chan eventReqResp {
	mux = http.NewServeMux()
	ts = httptest.NewServer(mux)
	control := make(chan eventReqResp)

	mux.HandleFunc("/v1/events", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "POST")
		body, _ := io.ReadAll(r.Body)
		respCh := make(chan int)

		select {
		case control <- eventReqResp{Body: body, Response: respCh}:
			status := <-respCh
			w.WriteHeader(status)
			if status == 201 {
				fmt.Fprint(w, `{"id":"event-id"}`)
			}
		default:
			w.WriteHeader(201)
			fmt.Fprint(w, `{"id":"event-id"}`)
		}
	})

	if DefaultClient.eventsWorker != nil {
		DefaultClient.eventsWorker.Stop()
	}

	config := newConfig(Configuration{APIKey: "badgers", Endpoint: ts.URL})
	*DefaultClient.Config = *config
	DefaultClient.eventsWorker = NewEventsWorker(config)

	return control
}

func teardown() {
	if DefaultClient.eventsWorker != nil {
		DefaultClient.eventsWorker.Stop()
	}
	*DefaultClient.Config = defaultConfig
	DefaultClient.beforeNotifyHandlers = nil
	DefaultClient.beforeEventHandlers = nil
}

func TestDefaultConfig(t *testing.T) {
	if Config.APIKey != "" {
		t.Errorf("Expected Config.APIKey to be empty by default. expected=%#v result=%#v", "", Config.APIKey)
	}
}

func TestConfigure(t *testing.T) {
	Configure(Configuration{APIKey: "badgers"})
	if Config.APIKey != "badgers" {
		t.Errorf("Expected Configure to override config.APIKey. expected=%#v actual=%#v", "badgers", Config.APIKey)
	}
}

func TestNotify(t *testing.T) {
	setup(t)
	defer teardown()

	res, _ := Notify(errors.New("Cobras!"))

	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}

	Flush()

	if !testRequestCount(t, 1) {
		return
	}

	testNoticePayload(t, requests[0].decodeJSON())
}

func TestNotifyWithContext(t *testing.T) {
	setup(t)
	defer teardown()

	context := Context{"foo": "bar"}
	Notify("Cobras!", context)
	Flush()

	if !testRequestCount(t, 1) {
		return
	}

	payload := requests[0].decodeJSON()
	if !testNoticePayload(t, payload) {
		return
	}

	assertContext(t, payload, context)
}

func TestNotifyWithErrorClass(t *testing.T) {
	setup(t)
	defer teardown()

	Notify("Cobras!", ErrorClass{"Badgers"})
	Flush()

	if !testRequestCount(t, 1) {
		return
	}

	payload := requests[0].decodeJSON()
	error_payload, _ := payload["error"].(map[string]any)
	sent_klass, _ := error_payload["class"].(string)

	if !testNoticePayload(t, payload) {
		return
	}

	if sent_klass != "Badgers" {
		t.Errorf("Custom error class should override default. expected=%v actual=%#v.", "Badgers", sent_klass)
		return
	}
}

func TestNotifyWithTags(t *testing.T) {
	setup(t)
	defer teardown()

	Notify("Cobras!", Tags{"timeout", "http"})
	Flush()

	if !testRequestCount(t, 1) {
		return
	}

	payload := requests[0].decodeJSON()
	error_payload, _ := payload["error"].(map[string]any)
	sent_tags, _ := error_payload["tags"].([]any)

	if !testNoticePayload(t, payload) {
		return
	}

	if got, want := sent_tags, []any{"timeout", "http"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Custom error class should override default. expected=%#v actual=%#v.", want, got)
		return
	}
}

func TestNotifyWithFingerprint(t *testing.T) {
	setup(t)
	defer teardown()

	Notify("Cobras!", Fingerprint{"Badgers"})
	Flush()

	if !testRequestCount(t, 1) {
		return
	}

	payload := requests[0].decodeJSON()
	error_payload, _ := payload["error"].(map[string]any)
	sent_fingerprint, _ := error_payload["fingerprint"].(string)

	if !testNoticePayload(t, payload) {
		return
	}

	if sent_fingerprint != "Badgers" {
		t.Errorf("Custom fingerprint should override default. expected=%v actual=%#v.", "Badgers", sent_fingerprint)
		return
	}
}

func TestNotifyWithRequest(t *testing.T) {
	setup(t)
	defer teardown()

	reqUrl := "/reqPath?qKey=qValue"
	var req *http.Request

	// Make sure nil request doesn't panic
	Notify("Cobras!", req)

	// Test a request with query data without form
	req = httptest.NewRequest("GET", reqUrl, nil)
	Notify("Cobras!", req)
	Flush()

	// Test a request with form and query data
	req = httptest.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Accept", "application/test-data")
	req.Form = url.Values{"fKey": {"fValue"}}
	Notify("Cobras!", req)
	Flush()

	if !testRequestCount(t, 3) {
		return
	}

	// Request[0] - Valid error means we properly handled a nil value
	if error := requests[0].decodeJSON()["error"]; error == nil {
		t.Errorf("Request error should be populated.")
	}

	// Request[1] - Checks URL & query extraction
	payload := requests[1].decodeJSON()
	request_payload, _ := payload["request"].(map[string]any)

	if url, _ := request_payload["url"].(string); url != reqUrl {
		t.Errorf("Request URL should be extracted. expected=%v actual=%#v.", "/fail", url)
		return
	}

	params, _ := request_payload["params"].(map[string]any)
	values, _ := params["qKey"].([]any)
	if len(params) != 1 || len(values) != 1 || values[0] != "qValue" {
		t.Errorf("Request params should be extracted. expected=%v actual=%#v.", req.Form, params)
	}

	// Request[2] - Checks header & form extraction
	payload = requests[2].decodeJSON()
	request_payload, _ = payload["request"].(map[string]any)

	if !testNoticePayload(t, payload) {
		return
	}

	cgi, _ := request_payload["cgi_data"].(map[string]any)
	if len(cgi) != 1 || cgi["HTTP_ACCEPT"] != "application/test-data" {
		t.Errorf("Request cgi_data should be extracted. expected=%v actual=%#v.", req.Header, cgi)
	}

	params, _ = request_payload["params"].(map[string]any)
	values, _ = params["fKey"].([]any)
	if len(params) != 1 || len(values) != 1 || values[0] != "fValue" {
		t.Errorf("Request params should be extracted. expected=%v actual=%#v.", req.Form, params)
	}
}

func TestMonitor(t *testing.T) {
	setup(t)
	defer teardown()

	defer func() {
		_ = recover()

		if !testRequestCount(t, 1) {
			return
		}

		testNoticePayload(t, requests[0].decodeJSON())
	}()

	defer Monitor()

	panic("Cobras!")
}

func TestNotifyWithHandler(t *testing.T) {
	setup(t)
	defer teardown()

	BeforeNotify(func(n *Notice) error {
		n.Fingerprint = "foo bar baz"
		return nil
	})
	Notify(errors.New("Cobras!"))
	Flush()

	payload := requests[0].decodeJSON()
	error_payload, _ := payload["error"].(map[string]any)
	sent_fingerprint, _ := error_payload["fingerprint"].(string)

	if !testRequestCount(t, 1) {
		return
	}

	if sent_fingerprint != "foo bar baz" {
		t.Errorf("Handler fingerprint should override default. expected=%v actual=%#v.", "foo bar baz", sent_fingerprint)
		return
	}
}

func TestNotifyWithHandlerError(t *testing.T) {
	setup(t)
	defer teardown()

	err := fmt.Errorf("Skipping this notification")

	BeforeNotify(func(n *Notice) error {
		return err
	})
	_, notifyErr := Notify(errors.New("Cobras!"))
	Flush()

	if !testRequestCount(t, 0) {
		return
	}

	if notifyErr != err {
		t.Errorf("Notify should return error from handler. expected=%v actual=%#v.", err, notifyErr)
		return
	}
}

// Helper functions.

func assertContext(t *testing.T, payload hash, expected Context) {
	var request, context hash
	var ok bool

	request, ok = payload["request"].(map[string]any)
	if !ok {
		t.Errorf("Missing request in payload actual=%#v.", payload)
		return
	}

	context, ok = request["context"].(map[string]any)
	if !ok {
		t.Errorf("Missing context in request payload actual=%#v.", request)
		return
	}

	for k, v := range expected {
		if context[k] != v {
			t.Errorf("Expected context to include hash. expected=%#v actual=%#v", expected, context)
			return
		}
	}
}

func testRequestCount(t *testing.T, num int) bool {
	if len(requests) != num {
		t.Errorf("Expected %v request to have been made. expected=%#v actual=%#v", num, num, len(requests))
		return false
	}
	return true
}

func testNoticePayload(t *testing.T, payload hash) bool {
	for _, key := range []string{"notifier", "error", "request", "server"} {
		switch payload[key].(type) {
		case map[string]any:
			// OK
		default:
			t.Errorf("Expected payload to include %v hash. expected=%#v actual=%#v", key, key, payload)
			return false
		}
	}
	return true
}

func TestEvent(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 1})

	eventData := map[string]any{
		"message": "test message",
		"user_id": 123,
	}

	err := Event("test_event", eventData)
	if err != nil {
		t.Errorf("Expected Event() to return no error. actual=%#v", err)
	}

	req := <-control
	req.Response <- 201

	lines := strings.Split(strings.TrimSpace(string(req.Body)), "\n")
	if len(lines) != 1 {
		t.Fatalf("Expected 1 JSONL event. actual=%d lines", len(lines))
	}

	var event map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &event); err != nil {
		t.Fatalf("Failed to parse JSONL event: %v", err)
	}
	if eventType, ok := event["event_type"].(string); !ok || eventType != "test_event" {
		t.Errorf("Expected event_type 'test_event'. actual=%#v", event["event_type"])
	}
	if message, ok := event["message"].(string); !ok || message != "test message" {
		t.Errorf("Expected message 'test message'. actual=%#v", event["message"])
	}
	if _, ok := event["ts"].(string); !ok {
		t.Errorf("Expected ts field to be present. actual=%#v", event)
	}
}

func TestEventBatching(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 2})

	Event("event1", map[string]any{"data": "first"})

	select {
	case <-control:
		t.Errorf("Expected no requests before batch size reached")
	case <-time.After(50 * time.Millisecond):
	}

	Event("event2", map[string]any{"data": "second"})

	select {
	case req := <-control:
		req.Response <- 201
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Expected 1 batch request when batch size reached")
	}
}

func TestEventTimeout(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 10, EventsTimeout: 50 * time.Millisecond})

	Event("event1", map[string]any{"data": "first"})

	select {
	case <-control:
		t.Errorf("Expected no immediate requests")
	case <-time.After(25 * time.Millisecond):
	}

	select {
	case req := <-control:
		req.Response <- 201
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Expected 1 request after timeout")
	}
}

func TestEventContextCancellation(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())
	Configure(Configuration{EventsBatchSize: 10, Context: ctx})

	Event("test_event", map[string]any{"data": "should be flushed"})

	select {
	case <-control:
		t.Errorf("Expected no requests before context cancellation")
	case <-time.After(50 * time.Millisecond):
	}

	cancel()

	select {
	case req := <-control:
		req.Response <- 201
		lines := strings.Split(strings.TrimSpace(string(req.Body)), "\n")
		var event map[string]any
		json.Unmarshal([]byte(lines[0]), &event)

		if event["data"] != "should be flushed" {
			t.Errorf("Expected event data 'should be flushed'. actual=%v", event["data"])
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Expected 1 request after context cancellation")
	}
}

func parseEvents(body string) []map[string]any {
	lines := strings.Split(strings.TrimSpace(body), "\n")
	events := make([]map[string]any, len(lines))

	for i, line := range lines {
		json.Unmarshal([]byte(line), &events[i])
	}

	return events
}

func TestEventMaxQueueSize(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 10, EventsMaxQueueSize: 2})

	Event("old_event", map[string]any{"data": "should be dropped"})
	Event("middle_event", map[string]any{"data": "middle"})
	Event("new_event", map[string]any{"data": "newest"})

	DefaultClient.eventsWorker.Flush()

	select {
	case req := <-control:
		req.Response <- 201
		events := parseEvents(string(req.Body))

		if len(events) != 2 {
			t.Fatalf("Expected 2 events. actual=%d", len(events))
		}

		expectedData := []string{"middle", "newest"}
		for i, expected := range expectedData {
			if events[i]["data"] != expected {
				t.Errorf("Expected event %d data '%s'. actual=%v", i, expected, events[i]["data"])
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Expected 1 request")
	}
}

func TestEventFailureRecovery(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{
		EventsBatchSize:    2,
		EventsMaxQueueSize: 5,
		EventsMaxRetries:   3,
		EventsTimeout:      50 * time.Millisecond,
	})

	Event("1", map[string]any{"data": "first"})
	Event("2", map[string]any{"data": "second"})

	req1 := <-control
	fmt.Println("First attempt", string(req1.Body))
	req1.Response <- 500

	req2 := <-control
	fmt.Println("Second attempt", string(req2.Body))
	req2.Response <- 500

	req3 := <-control
	fmt.Println("Third attempt", string(req3.Body))
	req3.Response <- 201
}

func TestEventQueueSizeIncludesPendingBatches(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{
		EventsBatchSize:    2,
		EventsMaxQueueSize: 3,
		EventsMaxRetries:   3,
	})

	Event("1", map[string]any{"data": "1"})
	Event("2", map[string]any{"data": "2"})

	req1 := <-control
	req1.Response <- 500

	Event("3", map[string]any{"data": "3"})
	Event("4", map[string]any{"data": "4"})

	time.Sleep(50 * time.Millisecond)

	DefaultClient.eventsWorker.Flush()

	req2 := <-control
	req2.Response <- 201

	req3 := <-control
	events := parseEvents(string(req3.Body))

	if len(events) != 1 {
		t.Errorf("Expected 1 event after dropping. actual=%d", len(events))
	}
	if events[0]["data"] != "4" {
		t.Errorf("Expected newest event '4'. actual=%v", events[0]["data"])
	}
	req3.Response <- 201
}

func TestEventThrottling(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{
		EventsBatchSize:    2,
		EventsMaxRetries:   2,
		EventsThrottleWait: 50 * time.Millisecond,
	})

	Event("1", map[string]any{"data": "1"})
	Event("2", map[string]any{"data": "2"})

	req1 := <-control
	req1.Response <- 429

	Event("3", map[string]any{"data": "3"})
	Event("4", map[string]any{"data": "4"})

	req2 := <-control
	req2.Response <- 429

	req3 := <-control
	req3.Response <- 201

	select {
	case <-control:
		t.Errorf("Expected no more batches, events 3 and 4 should have been dropped during throttling")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestEventMultipleBatchRetryOrdering(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{
		EventsBatchSize:  2,
		EventsMaxRetries: 3,
	})

	Event("1", map[string]any{"data": "1"})
	Event("2", map[string]any{"data": "2"})

	req1 := <-control
	req1.Response <- 500

	Event("3", map[string]any{"data": "3"})
	Event("4", map[string]any{"data": "4"})

	req2 := <-control
	events := parseEvents(string(req2.Body))
	if events[0]["data"] != "1" || events[1]["data"] != "2" {
		t.Errorf("Expected first batch to retry. actual=%v", events)
	}
	req2.Response <- 500

	req3 := <-control
	events = parseEvents(string(req3.Body))
	if events[0]["data"] != "1" || events[1]["data"] != "2" {
		t.Errorf("Expected first batch to retry again. actual=%v", events)
	}
	req3.Response <- 201

	req4 := <-control
	events = parseEvents(string(req4.Body))
	if len(events) != 2 {
		t.Fatalf("Expected second batch. actual=%d events", len(events))
	}
	if events[0]["data"] != "3" || events[1]["data"] != "4" {
		t.Errorf("Expected second batch after first succeeds. actual=%v", events)
	}
	req4.Response <- 201
}

func TestEventShutdownWithPendingRetries(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())
	Configure(Configuration{
		EventsBatchSize:  2,
		EventsMaxRetries: 3,
		Context:          ctx,
	})

	Event("1", map[string]any{"data": "1"})
	Event("2", map[string]any{"data": "2"})

	req1 := <-control
	req1.Response <- 500

	cancel()

	req2 := <-control
	events := parseEvents(string(req2.Body))
	if len(events) != 2 {
		t.Fatalf("Expected 2 events in retry batch on shutdown. actual=%d", len(events))
	}
	if events[0]["data"] != "1" || events[1]["data"] != "2" {
		t.Errorf("Expected failed batch to be flushed on shutdown. actual=%v", events)
	}
	req2.Response <- 201
}

func TestEventWithHandler(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 1})

	BeforeEvent(func(event map[string]any) error {
		event["modified"] = true
		return nil
	})

	Event("test_event", map[string]any{"data": "original"})

	req := <-control
	req.Response <- 201

	lines := strings.Split(strings.TrimSpace(string(req.Body)), "\n")
	var event map[string]any
	json.Unmarshal([]byte(lines[0]), &event)

	if event["modified"] != true {
		t.Errorf("Expected handler to modify event. actual=%v", event)
	}
	if event["data"] != "original" {
		t.Errorf("Expected original data to be preserved. actual=%v", event["data"])
	}
}

func TestEventWithHandlerError(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 1})

	err := fmt.Errorf("skip this event")
	BeforeEvent(func(event map[string]any) error {
		return err
	})

	eventErr := Event("test_event", map[string]any{"data": "test"})

	if eventErr != err {
		t.Errorf("Expected Event to return handler error. actual=%v", eventErr)
	}

	select {
	case <-control:
		t.Errorf("Expected no event to be sent when handler returns error")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestEventWithHandlerDropped(t *testing.T) {
	control := setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 1})

	BeforeEvent(func(event map[string]any) error {
		return ErrEventDropped
	})

	eventErr := Event("test_event", map[string]any{"data": "test"})

	if eventErr != nil {
		t.Errorf("Expected Event to return nil when ErrEventDropped. actual=%v", eventErr)
	}

	select {
	case <-control:
		t.Errorf("Expected no event to be sent when handler returns ErrEventDropped")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestHandlerCallsHandler(t *testing.T) {
	mockHandler := &MockedHandler{}
	mockHandler.On("ServeHTTP").Return()

	handler := Handler(mockHandler)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	mockHandler.AssertCalled(t, "ServeHTTP")
}
