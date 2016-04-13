package honeybadger

import "testing"

type TestLogger struct{}

func (l *TestLogger) Printf(format string, v ...interface{}) {}

type TestBackend struct{}

func (l *TestBackend) Notify(f Feature, p Payload) (err error) {
	return
}

func TestUpdateConfig(t *testing.T) {
	config := &Configuration{}
	logger := &TestLogger{}
	backend := &TestBackend{}
	config.update(&Configuration{
		Logger:  logger,
		Backend: backend,
		Root:    "/tmp/foo",
		Silent:  true,
	})

	if config.Logger != logger {
		t.Errorf("Expected config to update logger expected=%#v actual=%#v", logger, config.Logger)
	}
	if config.Backend != backend {
		t.Errorf("Expected config to update backend expected=%#v actual=%#v", backend, config.Backend)
	}
	if config.Root != "/tmp/foo" {
		t.Errorf("Expected config to update root expected=%#v actual=%#v", "/tmp/foo", config.Root)
	}
	if !config.Silent {
		t.Errorf("Expected config to update silent expected=%#v actual=%#v", true, config.Silent)
	}
}

func TestReplaceConfigPointer(t *testing.T) {
	config := Configuration{Root: "/tmp/foo"}
	root := &config.Root
	config = Configuration{Root: "/tmp/bar"}
	if *root != "/tmp/bar" {
		t.Errorf("Expected updated config to update pointer expected=%#v actual=%#v", "/tmp/bar", *root)
	}
}
