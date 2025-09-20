package honeybadger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	mux           *http.ServeMux
	ts            *httptest.Server
	requests      []*HTTPRequest
	eventRequests []*HTTPRequest
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
	body, _ := ioutil.ReadAll(r.Body)
	return &HTTPRequest{r, body}
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

func setupEvents(t *testing.T) {
	mux = http.NewServeMux()
	ts = httptest.NewServer(mux)
	eventRequests = []*HTTPRequest{}
	mux.HandleFunc("/v1/events",
		func(w http.ResponseWriter, r *http.Request) {
			assertMethod(t, r, "POST")
			eventRequests = append(eventRequests, newHTTPRequest(r))
			w.WriteHeader(201)
			fmt.Fprint(w, `{"id":"event-id"}`)
		},
	)

	*DefaultClient.Config = *newConfig(Configuration{APIKey: "badgers", Endpoint: ts.URL})
}

func teardown() {
	*DefaultClient.Config = defaultConfig
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
	error_payload, _ := payload["error"].(map[string]interface{})
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
	error_payload, _ := payload["error"].(map[string]interface{})
	sent_tags, _ := error_payload["tags"].([]interface{})

	if !testNoticePayload(t, payload) {
		return
	}

	if got, want := sent_tags, []interface{}{"timeout", "http"}; !reflect.DeepEqual(got, want) {
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
	error_payload, _ := payload["error"].(map[string]interface{})
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
	request_payload, _ := payload["request"].(map[string]interface{})

	if url, _ := request_payload["url"].(string); url != reqUrl {
		t.Errorf("Request URL should be extracted. expected=%v actual=%#v.", "/fail", url)
		return
	}

	params, _ := request_payload["params"].(map[string]interface{})
	values, _ := params["qKey"].([]interface{})
	if len(params) != 1 || len(values) != 1 || values[0] != "qValue" {
		t.Errorf("Request params should be extracted. expected=%v actual=%#v.", req.Form, params)
	}

	// Request[2] - Checks header & form extraction
	payload = requests[2].decodeJSON()
	request_payload, _ = payload["request"].(map[string]interface{})

	if !testNoticePayload(t, payload) {
		return
	}

	cgi, _ := request_payload["cgi_data"].(map[string]interface{})
	if len(cgi) != 1 || cgi["HTTP_ACCEPT"] != "application/test-data" {
		t.Errorf("Request cgi_data should be extracted. expected=%v actual=%#v.", req.Header, cgi)
	}

	params, _ = request_payload["params"].(map[string]interface{})
	values, _ = params["fKey"].([]interface{})
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
	error_payload, _ := payload["error"].(map[string]interface{})
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

	request, ok = payload["request"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing request in payload actual=%#v.", payload)
		return
	}

	context, ok = request["context"].(map[string]interface{})
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
		case map[string]interface{}:
			// OK
		default:
			t.Errorf("Expected payload to include %v hash. expected=%#v actual=%#v", key, key, payload)
			return false
		}
	}
	return true
}

func TestEvent(t *testing.T) {
	setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 1})

	eventData := map[string]interface{}{
		"message": "test message",
		"user_id": 123,
	}

	err := Event("test_event", eventData)
	if err != nil {
		t.Errorf("Expected Event() to return no error. actual=%#v", err)
	}

	if len(eventRequests) != 1 {
		t.Fatalf("Expected 1 event request. actual=%d", len(eventRequests))
	}

	payload := eventRequests[0].decodeJSON()
	events, ok := payload["events"].([]interface{})
	if !ok || len(events) != 1 {
		t.Fatalf("Expected batch format with 1 event. actual=%#v", payload)
	}

	event := events[0].(map[string]interface{})
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
	setupEvents(t)
	defer teardown()

	Configure(Configuration{EventsBatchSize: 2})

	Event("event1", map[string]interface{}{"data": "first"})

	if len(eventRequests) != 0 {
		t.Errorf("Expected no requests before batch size reached. actual=%d", len(eventRequests))
	}

	Event("event2", map[string]interface{}{"data": "second"})

	if len(eventRequests) != 1 {
		t.Fatalf("Expected 1 batch request when batch size reached. actual=%d", len(eventRequests))
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

func assertMethod(t *testing.T, r *http.Request, method string) {
	if r.Method != method {
		t.Errorf("Unexpected request method. actual=%#v expected=%#v", r.Method, method)
	}
}
