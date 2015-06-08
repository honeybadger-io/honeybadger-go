package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"testing"
)

func TestNotifyReturnsUUID(t *testing.T) {
	err := errors.New("Cobras!")
	var res string
	res = Notify(err)
	if uuid.Parse(res) == nil {
		t.Errorf("Expected honeybadger.Notify to be a UUID. result=%#v", res)
	}
}
