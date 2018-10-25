package honeybadger

import (
	"context"
	"testing"
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

	ctx := context.Background()

	ctx = client.SetContext(ctx, Context{"foo": "bar"})
	ctx = client.SetContext(ctx, Context{"bar": "baz"})

	var context Context
	if tmp, ok := ctx.Value(honeybadgerCtxKey).(*contextSync); ok {
		context = tmp.internal
	}

	if context == nil {
		t.Errorf("context value not placed in context.Context")
		t.FailNow()
	}

	if context["foo"] != "bar" {
		t.Errorf("Expected client to merge global context. expected=%#v actual=%#v", "bar", context["foo"])
	}

	if context["bar"] != "baz" {
		t.Errorf("Expected client to merge global context. expected=%#v actual=%#v", "baz", context["bar"])
	}
}
