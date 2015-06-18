package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	if Config.APIKey != "" {
		t.Errorf("Expected config.APIKey to be empty by default. expected=%#v result=%#v", "", Config.APIKey)
	}
}

func TestConfigure(t *testing.T) {
	Configure(Configuration{APIKey: "badgers"})
	if Config.APIKey != "badgers" {
		t.Errorf("Expected Configure to override config.APIKey. expected=%#v actual=%#v", "badgers", Config.APIKey)
	}
}

func TestNotify(t *testing.T) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)

	var requests []*http.Request
	mux.HandleFunc("/v1/notices",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			requests = append(requests, r)
			w.WriteHeader(201)
			fmt.Fprint(w, `{"id":"87ded4b4-63cc-480a-b50c-8abe1376d972"}`)
		},
	)

	Configure(Configuration{APIKey: "badgers", Endpoint: ts.URL})

	err := errors.New("Cobras!")
	var res string
	res = Notify(err)
	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}

	Flush()

	if len(requests) != 1 {
		t.Errorf("Expected 1 request to have been made. actual=%#v", len(requests))
	}
}

// Helper functions.

func testMethod(t *testing.T, r *http.Request, method string) {
	if r.Method != method {
		t.Errorf("Unexpected request method. actual=%#v expected=%#v", r.Method, method)
	}
}
