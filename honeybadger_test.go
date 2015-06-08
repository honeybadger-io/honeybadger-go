package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
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

func TestNotifyReturnsUUID(t *testing.T) {
	err := errors.New("Cobras!")
	var res string
	res = Notify(err)
	if uuid.Parse(res) == nil {
		t.Errorf("Expected Notify() to return a UUID. actual=%#v", res)
	}
}
