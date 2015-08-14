package honeybadger

import (
	"testing"
)

func TestNotifySilentMode(t *testing.T) {
	silent := true
	serv := server{Silent: &silent}
	if serv.Notify(Feature{}, nil) != nil {
		t.Error("Server in silent mode shouldn't send any errors to honeybadger")
	}
}
