package honeybadger

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

// The Logger interface is implemented by the standard log package and requires
// a limited subset of the interface implemented by log.Logger.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Configuration manages the configuration for the client.
type Configuration struct {
	APIKey             string
	Root               string
	Env                string
	Hostname           string
	Endpoint           string
	Sync               bool
	Timeout            time.Duration
	Logger             Logger
	Backend            Backend
	Context            context.Context
	EventsBatchSize      int
	EventsThrottleWait   time.Duration
	EventsTimeout        time.Duration
	EventsMaxQueueSize   int
	EventsMaxRetries     int
	EventsDropLogInterval time.Duration
}

func (c1 *Configuration) update(c2 *Configuration) *Configuration {
	if c2.APIKey != "" {
		c1.APIKey = c2.APIKey
	}
	if c2.Root != "" {
		c1.Root = c2.Root
	}
	if c2.Env != "" {
		c1.Env = c2.Env
	}
	if c2.Hostname != "" {
		c1.Hostname = c2.Hostname
	}
	if c2.Endpoint != "" {
		c1.Endpoint = c2.Endpoint
	}
	if c2.Timeout > 0 {
		c1.Timeout = c2.Timeout
	}
	if c2.Logger != nil {
		c1.Logger = c2.Logger
	}
	if c2.Backend != nil {
		c1.Backend = c2.Backend
	}
	if c2.Context != nil {
		c1.Context = c2.Context
	}
	if c2.EventsBatchSize > 0 {
		c1.EventsBatchSize = c2.EventsBatchSize
	}
	if c2.EventsTimeout > 0 {
		c1.EventsTimeout = c2.EventsTimeout
	}
	if c2.EventsMaxQueueSize > 0 {
		c1.EventsMaxQueueSize = c2.EventsMaxQueueSize
	}
	if c2.EventsMaxRetries > 0 {
		c1.EventsMaxRetries = c2.EventsMaxRetries
	}
	if c2.EventsThrottleWait > 0 {
		c1.EventsThrottleWait = c2.EventsThrottleWait
	}
	if c2.EventsDropLogInterval > 0 {
		c1.EventsDropLogInterval = c2.EventsDropLogInterval
	}

	c1.Sync = c2.Sync
	return c1
}

func newConfig(c Configuration) *Configuration {
	config := &Configuration{
		APIKey: GetEnv[string]("HONEYBADGER_API_KEY"),
		Root: GetEnv[string]("HONEYBADGER_ROOT", func() string {
			if val, err := os.Getwd(); err == nil {
				return val
			}
			return ""
		}),
		Env: GetEnv[string]("HONEYBADGER_ENV"),
		Hostname: GetEnv[string]("HONEYBADGER_HOSTNAME", func() string {
			if val, err := os.Hostname(); err == nil {
				return val
			}
			return ""
		}),
		Endpoint:           GetEnv[string]("HONEYBADGER_ENDPOINT", "https://api.honeybadger.io"),
		Timeout:            GetEnv[time.Duration]("HONEYBADGER_TIMEOUT", 3*time.Second),
		Logger:             log.New(os.Stderr, "[honeybadger] ", log.Flags()),
		Sync:               GetEnv[bool]("HONEYBADGER_SYNC", false),
		Context:            context.Background(),
		EventsThrottleWait:    GetEnv[time.Duration]("HONEYBADGER_EVENTS_THROTTLE_WAIT", 60*time.Second),
		EventsBatchSize:       GetEnv[int]("HONEYBADGER_EVENTS_BATCH_SIZE", 1000),
		EventsTimeout:         GetEnv[time.Duration]("HONEYBADGER_EVENTS_TIMEOUT", 30*time.Second),
		EventsMaxQueueSize:    GetEnv[int]("HONEYBADGER_EVENTS_MAX_QUEUE_SIZE", 100000),
		EventsMaxRetries:      GetEnv[int]("HONEYBADGER_EVENTS_MAX_RETRIES", 3),
		EventsDropLogInterval: GetEnv[time.Duration]("HONEYBADGER_EVENTS_DROP_LOG_INTERVAL", 60*time.Second),
	}
	config.update(&c)

	if config.Backend == nil {
		config.Backend = newServerBackend(config)
	}

	return config
}

func GetEnv[T any](key string, fallback ...any) T {
	val := os.Getenv(key)
	if val == "" {
		if len(fallback) > 0 {
			switch f := fallback[0].(type) {
			case func() T:
				return f()
			case T:
				return f
			}
		}
		var zero T
		return zero
	}

	switch any((*new(T))).(type) {
	case int:
		if v, err := strconv.Atoi(val); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return any(v).(T)
		}
	case bool:
		if v, err := strconv.ParseBool(val); err == nil {
			return any(v).(T)
		}
	case time.Duration:
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			return any(time.Duration(v)).(T)
		}
		if v, err := time.ParseDuration(val); err == nil {
			return any(v).(T)
		}
	case string:
		return any(val).(T)
	}

	if len(fallback) > 0 {
		switch f := fallback[0].(type) {
		case func() T:
			return f()
		case T:
			return f
		}
	}

	var zero T
	return zero
}
