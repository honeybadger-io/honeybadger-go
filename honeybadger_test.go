package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux           *http.ServeMux
	ts            *httptest.Server
	requests      []*HTTPRequest
	defaultConfig Configuration = *Config
)

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
			testMethod(t, r, "POST")
			requests = append(requests, newHTTPRequest(r))
			w.WriteHeader(201)
			fmt.Fprint(w, `{"id":"87ded4b4-63cc-480a-b50c-8abe1376d972"}`)
		},
	)

	*client.Config = *newConfig(Configuration{APIKey: "badgers", Endpoint: ts.URL})
}

func teardown() {
	*client.Config = defaultConfig
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

	res := Notify(errors.New("Cobras!"))

	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}

	Flush()

	testRequestCount(t, 1)
	testNoticePayload(t, requests[0].decodeJSON())
}

func TestNotifyWithContext(t *testing.T) {
	setup(t)
	defer teardown()

	context := Context{"foo": "bar"}
	Notify("Cobras!", context)
	Flush()

	testRequestCount(t, 1)

	payload := requests[0].decodeJSON()
	testNoticePayload(t, payload)
	assertContext(t, payload, context)
}

// Helper functions.

func assertContext(t *testing.T, payload hash, expected Context) {
	request, ok := payload["request"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing request in payload actual=%#v.", payload)
		return
	}

	context, ok := request["context"].(map[string]interface{})
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

func testRequestCount(t *testing.T, num int) {
	if len(requests) != num {
		t.Errorf("Expected %v request to have been made. expected=%#v actual=%#v", num, num, len(requests))
	}
}

func testNoticePayload(t *testing.T, payload hash) {
	for _, key := range []string{"notifier", "error", "request", "server"} {
		switch payload[key].(type) {
		case map[string]interface{}:
			// OK
		default:
			t.Errorf("Expected payload to include %v hash.", key)
		}
	}
}

func testMethod(t *testing.T, r *http.Request, method string) {
	if r.Method != method {
		t.Errorf("Unexpected request method. actual=%#v expected=%#v", r.Method, method)
	}
}
