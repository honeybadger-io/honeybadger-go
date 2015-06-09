# honeybadger-go

Go (golang) support for the :zap: [Honeybadger error
notifier](https://www.honeybadger.io/). Receive instant notification of panics
and errors in your Go applications.

## Installation

To install, grab the package from GitHub:

```sh
go get github.com/honeybadger-io/honeybadger-go
```

Then add an import to your application code:

```go
import "github.com/honeybadger-io/honeybadger-go"
```

Finally, configure your API key:

```go
honeybadger.Configure(honeybadger.Config{APIKey: "your api key"})
```

You can also configure Honeybadger via environment variables. See
[Configuration](./#configuration) for more information.

## Manually reporting panics

To report a panic manually, use `honeybadger.Notify`:

```go
if err != nil {
  honeybadger.Notify(err)
}
```

## Configuration

All available config options will eventually be listed here.
