package honeybadger

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

func (c Client) Notify(err interface{}) string {
	notice := newNotice(c.Config, newError(err, 1))
	go c.notify(notice)
	return notice.Token
}

func (c Client) notify(notice *Notice) {
	if err := c.Backend.Notify(Notices, notice); err != nil {
		panic(err)
	}
}

func NewClient(config Config) Client {
	defaultConfig := Config{
		APIKey:   getEnv("HONEYBADGER_API_KEY"),
		Env:      getEnv("HONEYBADGER_ENV"),
		Hostname: getHostname(),
		Endpoint: "https://api.honeybadger.io",
	}.merge(config)
	backend := Server{URL: &defaultConfig.Endpoint, APIKey: &defaultConfig.APIKey}
	return Client{
		Config:  &defaultConfig,
		Backend: backend,
	}
}
