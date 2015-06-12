package honeybadger

var (
	client  *Client
	Config  *Configuration
	Notices = Feature{"notices"}
)

type Feature struct {
	Endpoint string
}

func Configure(c Configuration) {
	client.Configure(c)
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

func Flush() {
	client.Flush()
}

func init() {
	client = New(Configuration{})
	Config = client.Config
}
