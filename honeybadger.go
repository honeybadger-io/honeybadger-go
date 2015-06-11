package honeybadger

var (
	client  Client
	config  *Config
	Notices = Feature{"notices"}
)

type Feature struct {
	Endpoint string
}

func Configure(c Config) {
	*client.Config = config.merge(c)
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
	client = NewClient(Config{})
	config = client.Config
}
