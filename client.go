package honeybadger

type Payload interface {
	toJSON() []byte
}

type Backend interface {
	Notify(feature Feature, payload Payload) error
}

type Client struct {
	Config  *Configuration
	Context *Context
	worker  worker
}

func (client *Client) Configure(config Configuration) {
	*client.Config = client.Config.merge(config)
}

func (client *Client) SetContext(context Context) {
	client.Context.Update(context)
}

func (c *Client) Flush() {
	c.worker.Flush()
}

func (c *Client) Notify(err interface{}, extra ...interface{}) string {
	extra = append([]interface{}{*c.Context}, extra...)
	notice := newNotice(c.Config, newError(err, 1), extra...)
	c.worker.Push(func() error {
		if err := c.Config.Backend.Notify(Notices, notice); err != nil {
			return err
		}
		return nil
	})
	return notice.Token
}

func New(c Configuration) *Client {
	config := newConfig(c)
	worker := newBufferedWorker(config)
	client := Client{
		Config:  config,
		worker:  worker,
		Context: &Context{},
	}

	return &client
}
