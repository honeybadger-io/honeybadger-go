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
honeybadger.Configure(honeybadger.Configuration{APIKey: "your api key"})
```

You can also configure Honeybadger via environment variables. See
[Configuration](#configuration) for more information.

## Manually reporting panics

To report a panic manually, use `honeybadger.Notify`:

```go
if err != nil {
  honeybadger.Notify(err)
}
```

## Creating a new client

In the same way that the log library provides a predefined "standard" logger,
honeybadger defines a standard client which may be access directly via
`honeybadger`. A new client may also be created by calling `honeybadger.New`:

```go
hb := honeybadger.New(honeybadger.Configuration{APIKey: "some other api key"})
hb.Notify("This error was reported by an alternate client.")
```

## Configuration

All available config options will eventually be listed here.

## Changelog

See https://github.com/honeybadger-io/honeybadger-go/releases

## Contributing

If you're adding a new feature, please [submit an issue](https://github.com/honeybadger-io/honeybadger-go/issues/new) as a preliminary step; that way you can be (moderately) sure that your pull request will be accepted.

### To contribute your code:

1. Fork it.
2. Create a topic branch `git checkout -b my_branch`
3. Commit your changes `git commit -am "Boom"`
3. Push to your branch `git push origin my_branch`
4. Send a [pull request](https://github.com/honeybadger-io/honeybadger-go/pulls)

### License

This library is MIT licensed. See the [LICENSE](https://raw.github.com/honeybadger-io/honeybadger-go/master/LICENSE) file in this repository for details.
