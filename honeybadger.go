package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
)

func Notify(err error) string {
	return uuid.NewRandom().String()
}
