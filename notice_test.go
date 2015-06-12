package honeybadger

import (
	"errors"
	"testing"
)

func TestNewNotice(t *testing.T) {
	err := errors.New("Cobras!")
	notice := newNotice(Config, newError(err, 0))
	if notice.ErrorMessage != "Cobras!" {
		t.Errorf("Unexpected value for notice.ErrorMessage. expected=%#v result=%#v", "Cobras!", notice.ErrorMessage)
	} else if notice.Error.err != err {
		t.Errorf("Unexpected value for notice.Error. expected=%#v result=%#v", err, notice.Error.err)
	}
}

func TestToJSON(t *testing.T) {
	err := errors.New("Cobras!")
	notice := newNotice(Config, newError(err, 0))
	notice.toJSON()
}
