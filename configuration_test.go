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
	new_config := config.merge(Configuration{
		Logger:  logger,
		Backend: backend,
		Root:    "/tmp/foo",
	})

	if config.Root != "" {
		t.Errorf("Merged config should not mutate original expected=%#v actual=%#v", "", config.Root)
	}

	if new_config.Logger != logger {
		t.Errorf("Expected config to merge logger expected=%#v actual=%#v", logger, new_config.Logger)
	}
	if new_config.Backend != backend {
		t.Errorf("Expected config to merge backend expected=%#v actual=%#v", backend, new_config.Backend)
	}
	if new_config.Root != "/tmp/foo" {
		t.Errorf("Expected config to merge root expected=%#v actual=%#v", "/tmp/foo", new_config.Root)
	}
}

func TestReplaceConfigPointer(t *testing.T) {
	config := Configuration{Root: "/tmp/foo"}
	root := &config.Root
	config = Configuration{Root: "/tmp/bar"}
	if *root != "/tmp/bar" {
		t.Errorf("Expected merged config to update pointer expected=%#v actual=%#v", "/tmp/bar", *root)
	}
}
