package honeybadger

import "testing"

type TestLogger struct{}

func (l *TestLogger) Printf(format string, v ...interface{}) {
}

func TestMergeConfig(t *testing.T) {
	config := Config{}
	logger := &TestLogger{}
	config.merge(Config{
		Logger: logger,
		Root:   "/tmp/foo",
	})
	if config.Logger != logger {
		t.Errorf("Expected config to merge logger expected=%#v actual=%#v", logger, config.Logger)
	}
	if config.Root != "/tmp/foo" {
		t.Errorf("Expected config to merge root expected=%#v actual=%#v", "/tmp/foo", config.Root)
	}
}
