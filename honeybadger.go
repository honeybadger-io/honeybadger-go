package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
)

type Config struct {
	APIKey string
}

var config Config

func Configure(c Config) {
	if c.APIKey != "" {
		config.APIKey = c.APIKey
	}
}

func Notify(err error) string {
	return uuid.NewRandom().String()
}

func init() {
	config = Config{
		APIKey: "",
	}
}
