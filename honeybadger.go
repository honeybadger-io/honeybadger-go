package honeybadger

import "os"

type Config struct {
	APIKey   string
	Env      string
	Hostname string
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

func getEnv(key string) string {
	return os.Getenv(key)
}

func getHostname() string {
	var hostname string
	hostname = getEnv("HONEYBADGER_HOSTNAME")
	if hostname == "" {
		if val, err := os.Hostname(); err == nil {
			hostname = val
		} else {
			panic(err)
		}
	}
	return hostname
}

func init() {
	config = Config{
		APIKey:   getEnv("HONEYBADGER_API_KEY"),
		Env:      getEnv("HONEYBADGER_ENV"),
		Hostname: getHostname(),
	}
}
