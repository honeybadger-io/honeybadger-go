package honeybadger

import "testing"

func TestNewClientConfig(t *testing.T) {
	client := NewClient(Config{APIKey: "lemmings"})
	if client.Config.APIKey != "lemmings" {
		t.Errorf("Expected NewClient to configure APIKey. expected=%#v actual=%#v", "lemmings", client.Config.APIKey)
	}
}

func TestConfigureClient(t *testing.T) {
	client := NewClient(Config{})
	client.Configure(Config{APIKey: "badgers"})
	if client.Config.APIKey != "badgers" {
		t.Errorf("Expected Configure to override config.APIKey. expected=%#v actual=%#v", "badgers", client.Config.APIKey)
	}
}
