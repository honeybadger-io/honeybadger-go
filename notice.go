package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
)

type hash map[string]interface{}

type Notice struct {
	APIKey       string
	Error        Error
	Token        string
	ErrorMessage string
	ErrorClass   string
	Hostname     string
	Env          string
	Backtrace    []*Frame
}

func (n *Notice) asJSON() *hash {
	return &hash{
		"api_key": n.APIKey,
		"notifier": &hash{
			"name":    "honeybadger",
			"url":     "https://github.com/honeybadger-io/honeybadger-go",
			"version": "0.0.0",
		},
		"error": &hash{
			"token":     n.Token,
			"message":   n.ErrorMessage,
			"class":     n.ErrorClass,
			"backtrace": n.Backtrace,
		},
		"server": &hash{
			"environment_name": n.Env,
			"hostname":         n.Hostname,
		},
	}
}

func (n *Notice) toJSON() []byte {
	if out, err := json.Marshal(n.asJSON()); err == nil {
		return out
	} else {
		panic(err)
	}
}

func newNotice(config *Configuration, err Error) *Notice {
	notice := Notice{
		APIKey:       config.APIKey,
		Error:        err,
		Token:        uuid.NewRandom().String(),
		ErrorMessage: err.Message,
		ErrorClass:   err.Class,
		Env:          config.Env,
		Hostname:     config.Hostname,
		Backtrace:    err.Stack,
	}

	return &notice
}
