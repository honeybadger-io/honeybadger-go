package honeybadger

import (
	"errors"
	"testing"
)

func TestNewNotice(t *testing.T) {
	err := errors.New("Cobras!")
	notice := newNotice(err)
	if notice.ErrorMessage != "Cobras!" {
		t.Errorf("Unexpected value for notice.ErrorMessage. expected=%#v result=%#v", "Cobras!", notice.ErrorMessage)
	} else if notice.Error != err {
		t.Errorf("Unexpected value for notice.Error. expected=%#v result=%#v", err, notice.Error)
	}
}
