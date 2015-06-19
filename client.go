package honeybadger

import "net/url"

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
	*client.Context = client.Context.merge(context)
}

func (c *Client) Flush() {
	c.worker.Flush()
}

func (c *Client) Notify(err interface{}, extra ...interface{}) string {
	notice := newNotice(c.Config, newError(err, 1))
	notice.setContext(*c.Context)
	for _, thing := range extra {
		switch thing := thing.(type) {
		case Context:
			notice.setContext(thing)
		case url.Values:
			notice.Params = thing
		}
	}
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
