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
	worker  Worker
}

func (c Client) Flush() {
	c.worker.Flush()
}

func (c Client) Notify(err interface{}) string {
	notice := newNotice(c.Config, newError(err, 1))
	c.worker.Push(func() error {
		if err := c.Backend.Notify(Notices, notice); err != nil {
			return err
		}
		return nil
	})
	return notice.Token
}

func NewClient(config Config) Client {
	defaultConfig := newConfig().merge(config)
	backend := Server{URL: &defaultConfig.Endpoint, APIKey: &defaultConfig.APIKey}
	client := Client{
		Config:  &defaultConfig,
		Backend: backend,
		worker:  newBufferedWorker(),
	}

	return client
}
