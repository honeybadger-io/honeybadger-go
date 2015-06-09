package honeybadger

import "os"

var (
	client  Client
	config  *Config
	Notices = Feature{"notices"}
)

type Config struct {
	APIKey   string
	Env      string
	Hostname string
	Endpoint string
	Backend  Backend
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

type Client struct {
	Config  *Config
	Backend Backend
}

func (c Client) Notify(err error) string {
	notice := newNotice(c.Config, err)
	if err := c.Backend.Notify(Notices, notice); err != nil {
		panic(err)
	}
	return notice.Token
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
	return client.Notify(err)
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

func NewClient(config Config) Client {
	defaultConfig := Config{
		APIKey:   getEnv("HONEYBADGER_API_KEY"),
		Env:      getEnv("HONEYBADGER_ENV"),
		Hostname: getHostname(),
		Endpoint: "https://api.honeybadger.io",
	}
	backend := Server{URL: &defaultConfig.Endpoint, APIKey: &defaultConfig.APIKey}
	return Client{
		Config:  &defaultConfig,
		Backend: backend,
	}
}

func init() {
	client = NewClient(Config{})
	config = client.Config
}
