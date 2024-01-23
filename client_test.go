package honeybadger

import (
	"context"
	"fmt"
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
		Config: newConfig(*backendConfig),
		worker: worker,
	}

	return client, worker, backend
}

func TestClientContext(t *testing.T) {
	backend := NewMemoryBackend()

	client := New(Configuration{
		APIKey:  "badgers",
		Backend: backend,
	})

	err := NewError(fmt.Errorf("which context is which"))

	hbCtx := Context{"user_id": 1}
	goCtx := Context{"request_id": "1234"}.WithContext(context.Background())

	_, nErr := client.Notify(err, hbCtx, goCtx)
	if nErr != nil {
		t.Fatal(nErr)
	}

	// Flush otherwise backend.Notices will be empty
	client.Flush()

	if len(backend.Notices) != 1 {
		t.Fatalf("Notices expected=%d actual=%d", 1, len(backend.Notices))
	}

	notice := backend.Notices[0]
	if notice.Context["user_id"] != 1 {
		t.Errorf("notice.Context[user_id] expected=%d actual=%v", 1, notice.Context["user_id"])
	}
	if notice.Context["request_id"] != "1234" {
		t.Errorf("notice.Context[request_id] expected=%q actual=%v", "1234", notice.Context["request_id"])
	}
}
