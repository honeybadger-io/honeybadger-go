package honeybadger

import "os"

var (
	client  Client
	config  *Config
	Notices = Feature{"notices"}
)

type Feature struct {
	Endpoint string
}

func Configure(c Config) {
	*client.Config = config.merge(c)
}

func Notify(err interface{}) string {
	return client.Notify(newError(err, 2))
}

func Monitor() {
	if err := recover(); err != nil {
		client.Notify(newError(err, 2))
		panic(err)
	}
}

func init() {
	client = NewClient(Config{})
	config = client.Config
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
