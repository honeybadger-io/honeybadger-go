package honeybadger

import "os"

type Config struct {
	APIKey   string
	Root     string
	Env      string
	Hostname string
	Endpoint string
}

func (c1 Config) merge(c2 Config) Config {
	if c2.APIKey != "" {
		c1.APIKey = c2.APIKey
	}
	if c2.Env != "" {
		c1.Env = c2.Env
	}
	if c2.Hostname != "" {
		c1.Hostname = c2.Hostname
	}
	if c2.Endpoint != "" {
		c1.Endpoint = c2.Endpoint
	}
	return c1
}

func newConfig() Config {
	return Config{
		APIKey:   getEnv("HONEYBADGER_API_KEY"),
		Root:     getPWD(),
		Env:      getEnv("HONEYBADGER_ENV"),
		Hostname: getHostname(),
		Endpoint: "https://api.honeybadger.io",
	}
}

// Private helper methods
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

func getPWD() string {
	var pwd string
	pwd = getEnv("HONEYBADGER_ROOT")
	if pwd == "" {
		if val, err := os.Getwd(); err == nil {
			pwd = val
		} else {
			panic(err)
		}
	}
	return pwd
}
