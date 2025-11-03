# slog-honeybadger

A Honeybadger handler for Go's slog package.
Send structured logs directly to Honeybadger as events.

## Features

- Structured logs as Honeybadger events
- Supports attributes, groups, and custom event types
- Chainable `With*` methods
- Log-level filtering with static or dynamic levels

## Install

```bash
go get github.com/honeybadger-io/honeybadger-go
```

**Requires Go 1.21+**

## Quick start

```go
import (
  "log/slog"
  "github.com/honeybadger-io/honeybadger-go"
  hbslog "github.com/honeybadger-io/honeybadger-go/slog"
)

func main() {
  client := honeybadger.New(honeybadger.Configuration{
    APIKey: "your-api-key",
  })
  logger := slog.New(hbslog.New(client))
  logger.Info("app started", "version", "1.0.0")
}
```

**Produces:**

```json
{
  "event_type": "log",
  "level": "INFO",
  "message": "app started",
  "version": "1.0.0"
}
```

## Event types

The default event type is `log`. You can set a custom event type for all
logs using `WithEventType`:

```go
audit := slog.New(hbslog.New(client).WithEventType("audit"))
audit.Info("user logged in", "user_id", 42)
```

**Produces:**

```json
{
  "event_type": "audit",
  "level": "INFO",
  "message": "user logged in",
  "user_id": 42
}
```

Set event type per log call with the `event_type` attribute:

```go
logger := slog.New(hbslog.New(client))

logger.Info("user signup", "event_type", "user_lifecycle", "user_id", 123)
logger.Info("payment processed", "event_type", "payment", "amount", 99.99)
logger.Info("regular log message", "foo", "bar")
```

**Produces three events with different types:**

```json
{"event_type": "user_lifecycle", "level": "INFO", "message": "user signup", "user_id": 123}
{"event_type": "payment", "level": "INFO", "message": "payment processed", "amount": 99.99}
{"event_type": "log", "level": "INFO", "message": "regular log message", "foo": "bar"}
```

## Attributes and groups

```go
handler := hbslog.New(client).
  WithAttrs([]slog.Attr{slog.String("service", "api")}).
  WithGroup("http")

logger := slog.New(handler)
logger.Info("request handled", "status", 200, "method", "POST")
```

**Produces:**

```json
{
  "event_type": "log",
  "level": "INFO",
  "message": "request handled",
  "service": "api",
  "http": {
    "status": 200,
    "method": "POST"
  }
}
```

## Log level filtering

Control which logs are sent to Honeybadger:

```go
// Only send WARN and above
handler := hbslog.New(client).WithLevel(slog.LevelWarn)
logger := slog.New(handler)

logger.Info("This is ignored")
logger.Warn("This is sent")
logger.Error("This is sent")
```

## Dynamic level changes

```go
levelVar := new(slog.LevelVar)
levelVar.Set(slog.LevelInfo)

handler := hbslog.New(client).WithLevel(levelVar)
logger := slog.New(handler)

// Change level at runtime
levelVar.Set(slog.LevelDebug) // Now debug logs will be sent
```

## License

MIT © Honeybadger.io
