package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
)

type hash map[string]interface{}

type Notice struct {
	Error           error
	Token           string
	ErrorMessage    string
	Hostname        string
	EnvironmentName string
}

func (n *Notice) asJSON() *hash {
	return &hash{
		"notifier": &hash{
			"name":    "honeybadger",
			"url":     "https://github.com/honeybadger-io/honeybadger-go",
			"version": "0.0.0",
		},
		"error": &hash{
			"token":     n.Token,
			"message":   n.ErrorMessage,
			"backtrace": []map[string]interface{}{},
		},
		"server": &hash{
			"environment_name": n.EnvironmentName,
			"hostname":         n.Hostname,
		},
	}
}

func (n *Notice) toJSON() string {
	if out, err := json.Marshal(n.asJSON()); err == nil {
		return string(out)
	} else {
		panic(err)
	}
}

func newNotice(err error) *Notice {
	notice := Notice{
		Error:           err,
		Token:           uuid.NewRandom().String(),
		ErrorMessage:    err.Error(),
		EnvironmentName: "production",
		Hostname:        "localhost",
	}

	return &notice
}
