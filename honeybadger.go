package honeybadger

import "os"

type Config struct {
	APIKey          string
	EnvironmentName string
	Hostname        string
}

var config Config

func Configure(c Config) {
	if c.APIKey != "" {
		config.APIKey = c.APIKey
	}
}

func Notify(err error) string {
	notice := newNotice(&config, err)
	return notice.Token
}

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	config = Config{
		APIKey:          "",
		EnvironmentName: "",
		Hostname:        hostname,
	}
}
