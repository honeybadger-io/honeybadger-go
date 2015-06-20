package honeybadger

type Payload interface {
	toJSON() []byte
}

type Backend interface {
	Notify(feature Feature, payload Payload) error
}

type Client struct {
	Config  *Configuration
	context *Context
	worker  worker
}

func (client *Client) Configure(config Configuration) {
	*client.Config = client.Config.merge(config)
}

func (client *Client) SetContext(context Context) {
	client.context.Update(context)
}

func (c *Client) Flush() {
	c.worker.Flush()
}

func (c *Client) Notify(err interface{}, extra ...interface{}) string {
	extra = append([]interface{}{*c.context}, extra...)
	notice := newNotice(c.Config, newError(err, 1), extra...)
	c.worker.Push(func() error {
		if err := c.Config.Backend.Notify(Notices, notice); err != nil {
			return err
		}
		return nil
	})
	return notice.Token
}

func (c *Client) Monitor() {
	if err := recover(); err != nil {
		client.Notify(newError(err, 1))
		panic(err)
	}
}

func New(c Configuration) *Client {
	config := newConfig(c)
	worker := newBufferedWorker(config)
	client := Client{
		Config:  config,
		worker:  worker,
		context: &Context{},
	}

	return &client
}
