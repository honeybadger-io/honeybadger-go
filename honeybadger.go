package honeybadger

type Config struct {
	APIKey string
}

var config Config

func Configure(c Config) {
	if c.APIKey != "" {
		config.APIKey = c.APIKey
	}
}

func Notify(err error) string {
	notice := newNotice(&config, err)
	return notice.Token
}

func init() {
	config = Config{
		APIKey: "",
	}
}
