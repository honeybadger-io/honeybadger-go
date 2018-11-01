package honeybadger

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const honeybadgerCtxKey = "honeybadger-go-ctx"

var noCastErr = fmt.Errorf("unable to cast value from context properly")

// The Payload interface is implemented by any type which can be handled by the
// Backend interface.
type Payload interface {
	toJSON() []byte
}

// The Backend interface is implemented by the server type by default, but a
// custom implementation may be configured by the user.
type Backend interface {
	Notify(feature Feature, payload Payload) error
}

type noticeHandler func(*Notice) error

// Client is the manager for interacting with the Honeybadger service. It holds
// the configuration and implements the public API.
type Client struct {
	Config               *Configuration
	worker               worker
	beforeNotifyHandlers []noticeHandler
}

// Configure updates the client configuration with the supplied config.
func (client *Client) Configure(config Configuration) {
	client.Config.update(&config)
}

// SetContext updates the client context with supplied context.
func (client *Client) SetContext(ctx context.Context, val Context) context.Context {
	var clientContext *contextSync

	tmp := ctx.Value(honeybadgerCtxKey)
	if tmp == nil {
		clientContext = &contextSync{
			internal: Context{},
		}
	} else {
		if cc, ok := tmp.(*contextSync); ok {
			clientContext = cc
		} else {
			panic(noCastErr)
		}
	}

	clientContext.Update(val)

	return context.WithValue(ctx, honeybadgerCtxKey, clientContext)
}

// Flush blocks until the worker has processed its queue.
func (client *Client) Flush() {
	client.worker.Flush()
}

// BeforeNotify adds a callback function which is run before a notice is
// reported to Honeybadger. If any function returns an error the notification
// will be skipped, otherwise it will be sent.
func (client *Client) BeforeNotify(handler func(notice *Notice) error) {
	client.beforeNotifyHandlers = append(client.beforeNotifyHandlers, handler)
}

// Notify reports the error err to the Honeybadger service.
func (client *Client) Notify(ctx context.Context, err interface{}, extra ...interface{}) (string, error) {
	val, ok := ctx.Value(honeybadgerCtxKey).(*contextSync)
	if ok {
		val.Lock()
		extra = append([]interface{}{val.internal}, extra...)
		val.Unlock()
	}

	notice := newNotice(client.Config, newError(err, 2), extra...)
	for _, handler := range client.beforeNotifyHandlers {
		if err := handler(notice); err != nil {
			return "", err
		}
	}
	workerErr := client.worker.Push(func() error {
		if err := client.Config.Backend.Notify(Notices, notice); err != nil {
			return err
		}
		return nil
	})
	if workerErr != nil {
		client.Config.Logger.Printf("worker error: %v\n", workerErr)
		return "", workerErr
	}
	return notice.Token, nil
}

// Monitor automatically reports panics which occur in the function it's called
// from. Must be deferred.
func (client *Client) Monitor(ctx context.Context) {
	if err := recover(); err != nil {
		client.Notify(ctx, newError(err, 2))
		client.Flush()
		panic(err)
	}
}

// Handler returns an http.Handler function which automatically reports panics
// to Honeybadger and then re-panics.
func (client *Client) Handler(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = client.SetContext(ctx, map[string]interface{}{
			"path":           r.URL.String(),
			"host":           r.Host,
			"header":         r.Header,
			"user-agent":     r.UserAgent(),
			"content-length": r.ContentLength,
			"method":         r.Method,
			"proto":          r.Proto,
			"remote-address": r.RemoteAddr,
			"request-uri":    r.RequestURI,
		})
		req := r.WithContext(ctx)

		defer func() {
			if err := recover(); err != nil {
				client.Notify(req.Context(), newError(err, 2), Params(r.Form), getCGIData(r), *r.URL)
				panic(err)
			}
		}()
		h.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// New returns a new instance of Client.
func New(c Configuration) *Client {
	config := newConfig(c)
	worker := newBufferedWorker(config)

	client := Client{
		Config: config,
		worker: worker,
	}

	return &client
}

func getCGIData(request *http.Request) CGIData {
	cgiData := CGIData{}
	replacer := strings.NewReplacer("-", "_")
	for k, v := range request.Header {
		key := "HTTP_" + replacer.Replace(strings.ToUpper(k))
		cgiData[key] = v[0]
	}
	return cgiData
}
