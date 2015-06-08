package honeybadger

import "code.google.com/p/go-uuid/uuid"

type Notice struct {
	Error           error
	Token           string
	ErrorMessage    string
	Hostname        string
	EnvironmentName string
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
