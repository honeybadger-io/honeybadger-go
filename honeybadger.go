package honeybadger

import (
	"net/http"
	"net/url"
	"strings"
)

var (
	// The global client.
	client *Client = New(Configuration{})

	// The global configuration (available through the client).
	Config *Configuration = client.Config

	// Notices is the feature for sending error reports.
	Notices = Feature{"notices"}
)

// A feature is provided by the API service. Its Endpoint maps to the
// collection endpoint of the /v1 API.
type Feature struct {
	Endpoint string
}

// CGI variables such as HTTP_METHOD.
type CGIData hash

// Request parameters.
type Params url.Values

// Configures the global client.
func Configure(c Configuration) {
	client.Configure(c)
}

// Set/merge the global context.
func SetContext(c Context) {
	client.SetContext(c)
}

// Notify reports the error err to the Honeybadger service.
//
// The first argument err may be an error, a string, or any other type in which
// case its formatted value will be used.
//
// It returns a string UUID which can be used to reference the error from the
// Honeybadger service.
func Notify(err interface{}, extra ...interface{}) string {
	return client.Notify(newError(err, 2), extra...)
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

// Returns an http.Handler function which automatically reports panics to
// Honeybadger including request data.
func Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				client.Notify(newError(err, 3), Params(r.Form), getCGIData(r), *r.URL)
				panic(err)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getCGIData(request *http.Request) CGIData {
	cgi_data := CGIData{}
	replacer := strings.NewReplacer("-", "_")
	for k, v := range request.Header {
		key := "HTTP_" + replacer.Replace(strings.ToUpper(k))
		cgi_data[key] = v[0]
	}
	return cgi_data
}
