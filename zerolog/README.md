# Zerolog adapter for Honeybadger Insights

Sends zerolog JSON logs to Honeybadger Insights as events.

Typical usage:

```go
import (
    "github.com/rs/zerolog"
    "github.com/honeybadger-io/honeybadger-go"
    hbzerolog "github.com/honeybadger-io/honeybadger-go/zerolog"
)

hb := honeybadger.New(honeybadger.Configuration{})
writer := hbzerolog.New(
    hb,
    hbzerolog.WithEventType("app_log"),
    hbzerolog.WithKeys("timestamp", "severity"),
)
log := zerolog.New(writer).With().Timestamp().Logger()
log.Info().Msg("hello")
```

## Options

### WithEventType(string)
Sets the default event type for all logs (default: "log"). Override
per-log by including an `event_type` field.

### WithKeys(timeKey, levelKey string)
Customize field names if your zerolog uses non-standard keys
(defaults: "time", "level"). The writer remaps the time field to
"ts" for Honeybadger.
