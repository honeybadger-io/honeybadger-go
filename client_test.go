package honeybadger

import (
	"sync"
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

	client := New(Configuration{})
	newContext := Context{"foo": "bar"}

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

func TestNotifyPushesTheEnvelope(t *testing.T) {
	client, worker, _ := mockClient(Configuration{})

	client.Notify("test")

	if worker.receivedEnvelope == false {
		t.Errorf("Expected client to push envelope")
	}
}

func TestNotifySyncMode(t *testing.T) {
	client, worker, backend := mockClient(Configuration{Sync: true})

	token, _ := client.Notify("test")

	if worker.receivedEnvelope == true {
		t.Errorf("Expected client to not push envelope")
	}
	if backend.notice.Token != token {
		t.Errorf("Notice should have been called on backend")
	}
}

type mockWorker struct {
	receivedEnvelope bool
}

func (w *mockWorker) Push(work envelope) error {
	w.receivedEnvelope = true
	return nil
}

func (w *mockWorker) Flush() {}

type mockBackend struct {
	notice *Notice
}

func (b *mockBackend) Notify(_ Feature, n Payload) error {
	b.notice = n.(*Notice)
	return nil
}

func mockClient(c Configuration) (Client, *mockWorker, *mockBackend) {
	worker := &mockWorker{}
	backend := &mockBackend{}
	backendConfig := &Configuration{Backend: backend}
	backendConfig.update(&c)

	client := Client{
		Config:  newConfig(*backendConfig),
		worker:  worker,
		context: newContextSync(),
	}

	return client, worker, backend
}
