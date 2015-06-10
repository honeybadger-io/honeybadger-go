package honeybadger

type Config struct {
	APIKey   string
	Env      string
	Hostname string
	Endpoint string
	Backend  Backend
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
