package honeybadger

type Payload interface {
	toJSON() []byte
}

type Backend interface {
	Notify(feature Feature, payload Payload) error
}

type Client struct {
	Config *Config
	worker Worker
}

func (c Client) Flush() {
	c.worker.Flush()
}

func (c Client) Notify(err interface{}) string {
	notice := newNotice(c.Config, newError(err, 1))
	c.worker.Push(func() error {
		if err := c.Config.Backend.Notify(Notices, notice); err != nil {
			return err
		}
		return nil
	})
	return notice.Token
}

func NewClient(c Config) *Client {
	config := newConfig(c)
	worker := newBufferedWorker(config)
	client := Client{
		Config: config,
		worker: worker,
	}

	return &client
}
