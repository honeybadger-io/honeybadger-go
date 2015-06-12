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

func TestNotifyReturnsUUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, "{\"id\":\"87ded4b4-63cc-480a-b50c-8abe1376d972\"}")
	}))
	defer ts.Close()
	Configure(Configuration{APIKey: "badgers", Endpoint: ts.URL})

	err := errors.New("Cobras!")
	var res string
	res = Notify(err)
	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}
}
