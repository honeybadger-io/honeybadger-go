package honeybadger

import "os"

var (
	config          Config
	BackendInstance Backend
	Notices         = Feature{"notices"}
)

type Config struct {
	APIKey   string
	Env      string
	Hostname string
	Endpoint string
}

type Feature struct {
	Endpoint string
}

type Payload interface {
	toJSON() []byte
}

type Backend interface {
	Notify(feature Feature, payload Payload) error
}

func Configure(c Config) {
	if c.APIKey != "" {
		config.APIKey = c.APIKey
	}
	if c.Env != "" {
		config.Env = c.Env
	}
	if c.Hostname != "" {
		config.Hostname = c.Hostname
	}
	if c.Endpoint != "" {
		config.Endpoint = c.Endpoint
	}
}

func Notify(err error) string {
	notice := newNotice(&config, err)
	if err := BackendInstance.Notify(Notices, notice); err != nil {
		panic(err)
	}
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
		Endpoint: "https://api.honeybadger.io",
	}
	BackendInstance = Server{URL: &config.Endpoint, APIKey: &config.APIKey}
}
