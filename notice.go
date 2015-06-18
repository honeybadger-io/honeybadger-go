package honeybadger

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"regexp"
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
	ProjectRoot  string
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
			"project_root":     n.ProjectRoot,
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

func composeStack(stack []*Frame, root string) (frames []*Frame) {
	if root == "" {
		return stack
	}

	re, err := regexp.Compile("^" + regexp.QuoteMeta(root))
	if err != nil {
		return stack
	}

	for _, frame := range stack {
		file := re.ReplaceAllString(frame.File, "[PROJECT_ROOT]")
		frames = append(frames, &Frame{
			File:   file,
			Number: frame.Number,
			Method: frame.Method,
		})
	}
	return
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
		Backtrace:    composeStack(err.Stack, config.Root),
		ProjectRoot:  config.Root,
	}

	return &notice
}
