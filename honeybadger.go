package honeybadger

var (
	// The global client.
	client *Client

	// The global configuration (available through the client).
	Config *Configuration

	// Notices is the feature for sending error reports.
	Notices = Feature{"notices"}
)

// A feature is provided by the API service. Its Endpoint maps to the
// collection endpoint of the /v1 API.
type Feature struct {
	Endpoint string
}

// Configures the global client.
func Configure(c Configuration) {
	client.Configure(c)
}

// Notify reports the error err to the Honeybadger service.
//
// The first argument err may be an error, a string, or any other type in which
// case its formatted value will be used.
//
// It returns a string UUID which can be used to reference the error from the
// Honeybadger service.
func Notify(err interface{}) string {
	return client.Notify(newError(err, 2))
}

// Monitor is used to automatically notify Honeybadger service of panics which
// happen inside the current function. In order to monitor for panics, defer a
// call to Monitor. For example:
// 	func main {
// 		defer honeybadger.Monitor()
// 		// Do risky stuff...
// 	}
// The Monitor function re-panics after the notification has been sent, so it's
// still up to the user to recover from panics if desired.
func Monitor() {
	if err := recover(); err != nil {
		client.Notify(newError(err, 2))
		panic(err)
	}
}

// Flush blocks until all data (normally sent in the background) has been sent
// to the Honeybadger service.
func Flush() {
	client.Flush()
}

func init() {
	client = New(Configuration{})
	Config = client.Config
}
