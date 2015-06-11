package honeybadger

import "testing"

type TestLogger struct{}

func (l *TestLogger) Printf(format string, v ...interface{}) {
}

func TestMergeConfig(t *testing.T) {
	config := Config{}
	logger := &TestLogger{}
	config = config.merge(Config{Logger: logger})
	if config.Logger != logger {
		t.Errorf("Expected config to merge logger expected=%#v actual=%#v", logger, config.Logger)
	}
}
