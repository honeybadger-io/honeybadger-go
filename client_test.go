package honeybadger

import (
	"testing"
	"sync"
)

func TestNewConfig(t *testing.T) {
	client := New(Configuration{APIKey: "lemmings"})
	if client.Config.APIKey != "lemmings" {
		t.Errorf("Expected New to configure APIKey. expected=%#v actual=%#v", "lemmings", client.Config.APIKey)
	}
}

func TestConfigureClient(t *testing.T) {
	client := New(Configuration{})
	client.Configure(Configuration{APIKey: "badgers"})
	if client.Config.APIKey != "badgers" {
		t.Errorf("Expected Configure to override config.APIKey. expected=%#v actual=%#v", "badgers", client.Config.APIKey)
	}
}

func TestConfigureClientEndpoint(t *testing.T) {
	client := New(Configuration{})
	backend := client.Config.Backend.(*server)
	client.Configure(Configuration{Endpoint: "http://localhost:3000"})
	if *backend.URL != "http://localhost:3000" {
		t.Errorf("Expected Configure to update backend. expected=%#v actual=%#v", "http://localhost:3000", backend.URL)
	}
}

func TestClientContext(t *testing.T) {
	client := New(Configuration{})

	client.SetContext(Context{"foo": "bar"})
	client.SetContext(Context{"bar": "baz"})

	context := client.context.internal

	if context["foo"] != "bar" {
		t.Errorf("Expected client to merge global context. expected=%#v actual=%#v", "bar", context["foo"])
	}

	if context["bar"] != "baz" {
		t.Errorf("Expected client to merge global context. expected=%#v actual=%#v", "baz", context["bar"])
	}
}

func TestClientConcurrentContext(t *testing.T) {
	var wg sync.WaitGroup

	client     := New(Configuration{})
	newContext := Context{"foo":"bar"}

	wg.Add(2)

	go updateContext(&wg, client, newContext)
	go updateContext(&wg, client, newContext)

	wg.Wait()

	context := client.context.internal

	if context["foo"] != "bar" {
		t.Errorf("Expected context value. expected=%#v result=%#v", "bar", context["foo"])
	}
}

func updateContext(wg *sync.WaitGroup, client *Client, context Context) {
	client.SetContext(context)
	wg.Done()
}
