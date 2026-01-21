# Honeybadger for Go

[![Test](https://github.com/honeybadger-io/honeybadger-go/actions/workflows/tests.yml/badge.svg)](https://github.com/honeybadger-io/honeybadger-go/actions/workflows/tests.yml)

Go (golang) support for the :zap: [Honeybadger Exception Notifier](https://www.honeybadger.io/). Receive instant notification of panics and errors in your Go applications.

## Documentation and Support

For comprehensive documentation and support, [check out our documentation site](https://docs.honeybadger.io/lib/go/).

## Supported Go Versions

This library supports the last two major Go releases, consistent with the Go team's [release policy](https://go.dev/doc/devel/release):

- Go 1.25.x
- Go 1.24.x

Older versions may work but are not officially supported or tested.

## Changelog

See https://github.com/honeybadger-io/honeybadger-go/blob/master/CHANGELOG.md

## Development

Pull requests are welcome. If you're adding a new feature, please [submit an issue](https://github.com/honeybadger-io/honeybadger-go/issues/new) as a preliminary step; that way you can be (moderately) sure that your pull request will be accepted.

### To contribute your code:

1. Fork it.
2. Create a topic branch `git checkout -b my_branch`
3. Commit your changes `git commit -am "Boom"`
4. Push to your branch `git push origin my_branch`
5. Send a [pull request](https://github.com/honeybadger-io/honeybadger-go/pulls)

### Releasing

Releases are automated using [release-please](https://github.com/google-github-actions/release-please-action). When a PR is merged to master, a release PR is created (or updated) with version bumps and changelog entries based on [conventional commits](https://www.conventionalcommits.org/). Merging the release PR creates the GitHub release and tag.

### License

This library is MIT licensed. See the [LICENSE](https://raw.github.com/honeybadger-io/honeybadger-go/master/LICENSE) file in this repository for details.
