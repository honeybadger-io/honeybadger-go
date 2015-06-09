package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"reflect"
	"runtime"
	"strconv"
)

const MaxFrames = 20

type Frame struct {
	Number string `json:"number"`
	File   string `json:"file"`
	Method string `json:"method"`
}

type hash map[string]interface{}

type Notice struct {
	APIKey       string
	Error        error
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

func generateStack() (frames []*Frame) {
	stack := make([]uintptr, MaxFrames)
	length := runtime.Callers(5, stack[:])
	for _, pc := range stack[:length] {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}
		file, line := f.FileLine(pc)
		frame := &Frame{
			File:   file,
			Number: strconv.Itoa(line),
			Method: f.Name(),
		}
		frames = append(frames, frame)
	}

	return
}

func newNotice(config *Config, err error) *Notice {
	notice := Notice{
		APIKey:       config.APIKey,
		Error:        err,
		Token:        uuid.NewRandom().String(),
		ErrorMessage: err.Error(),
		ErrorClass:   reflect.TypeOf(err).String(),
		Env:          config.Env,
		Hostname:     config.Hostname,
		Backtrace:    generateStack(),
	}

	return &notice
}
