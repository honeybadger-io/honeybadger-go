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
	if config.APIKey != "" {
		t.Errorf("Expected config.APIKey to be empty by default. expected=%#v result=%#v", "", config.APIKey)
	}
}

func TestConfigure(t *testing.T) {
	Configure(Config{APIKey: "badgers"})
	if config.APIKey != "badgers" {
		t.Errorf("Expected Configure to override config.APIKey. expected=%#v actual=%#v", "badgers", config.APIKey)
	}
}

func TestClientConfig(t *testing.T) {
	Configure(Config{APIKey: "badgers"})
	if client.Config != config {
		t.Errorf("Expected client configuration to match global config. expected=%#v actual=%#v", config, client.Config)
	}
}

func TestNewClientConfig(t *testing.T) {
	client := NewClient(Config{APIKey: "lemmings"})
	if client.Config.APIKey != "lemmings" {
		t.Errorf("Expected NewClient to configure APIKey. expected=%#v actual=%#v", "lemmings", client.Config.APIKey)
	}
}

func TestNotifyReturnsUUID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, "{\"id\":\"87ded4b4-63cc-480a-b50c-8abe1376d972\"}")
	}))
	defer ts.Close()
	APIKey := "badgers"
	client.Backend = Server{APIKey: &APIKey, URL: &ts.URL}

	err := errors.New("Cobras!")
	var res string
	res = Notify(err)
	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}
}
