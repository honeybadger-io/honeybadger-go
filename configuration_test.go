package honeybadger

import "testing"

type TestLogger struct{}

func (l *TestLogger) Printf(format string, v ...interface{}) {}

type TestBackend struct{}

func (l *TestBackend) Notify(f Feature, p Payload) (err error) {
	return
}

func TestMergeConfig(t *testing.T) {
	config := Configuration{}
	logger := &TestLogger{}
	backend := &TestBackend{}
	config = config.merge(Configuration{
		Logger:  logger,
		Backend: backend,
		Root:    "/tmp/foo",
	})
	if config.Logger != logger {
		t.Errorf("Expected config to merge logger expected=%#v actual=%#v", logger, config.Logger)
	}
	if config.Backend != backend {
		t.Errorf("Expected config to merge backend expected=%#v actual=%#v", backend, config.Backend)
	}
	if config.Root != "/tmp/foo" {
		t.Errorf("Expected config to merge root expected=%#v actual=%#v", "/tmp/foo", config.Root)
	}
}
