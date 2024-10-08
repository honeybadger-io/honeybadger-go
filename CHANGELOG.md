# Change Log

All notable changes to this project will be documented in this file. See [Keep a
CHANGELOG](http://keepachangelog.com/) for how to update this file. This project
adheres to [Semantic Versioning](http://semver.org/).

## [0.8.0](https://github.com/honeybadger-io/honeybadger-go/compare/v0.7.0...v0.8.0) (2024-09-16)


### Features

* allow Notify to receive *http.Request ([#54](https://github.com/honeybadger-io/honeybadger-go/issues/54)) ([b0d2af0](https://github.com/honeybadger-io/honeybadger-go/commit/b0d2af07a031075d86b497c882f2649f25ca1404))

## [0.7.0](https://github.com/honeybadger-io/honeybadger-go/compare/v0.6.1...v0.7.0) (2024-04-25)


### Features

* allow errors to provide their own stack traces ([b8d3e83](https://github.com/honeybadger-io/honeybadger-go/commit/b8d3e83e6a36f7dac1b72e3ff7d1bf9cde4382da))

## [Unreleased][unreleased]

## [0.6.1] - 2024-04-12

### Fixed

- Implement error chains Unwrap method

## [0.6.0] - 2024-01-19

### Changed

- Updated dependencies

## [0.5.0] - 2019-10-17

### Added

- Added Sync mode

## [0.4.0] - 2018-07-18

### Added

- Ability to tag errors. -@izumin5210

## [0.3.0] - 2018-07-03

### Changed

- Remove deprecated metrics methods.

### Fixed

- Fixed concurrent map writes bug when calling `honeybadger.SetContext` from
  concurrent goroutines.

## [0.2.1] - 2017-09-14

### Fixed

- Previously, if you put `honeybadger.Monitor()` in your main func, the app
  could finish and exit before the error was sent to honeybadger. We now Flush
  notices before re-panicking.

## [0.2.0] - 2016-10-14

### Changed

- Sunset performance metrics. See
  http://blog.honeybadger.io/sunsetting-performance-metrics/

## [0.1.0] - 2016-05-12

### Added

- Use `honeybadger.MetricsHandler` to send us request metrics!

## [0.0.3] - 2016-04-13

### Added

- `honeybadger.NewNullBackend()`: creates a backend which swallows all errors
  and does not send them to Honeybadger. This is useful for development and
  testing to disable sending unnecessary errors. -@gaffneyc
- Tested against Go 1.5 and 1.6. -@gaffneyc

### Fixed

- Export Fingerprint fields. -@smeriwether
- Fix HB due to changes in shirou/gopsutil. -@kostyantyn

## [0.0.2] - 2016-03-28

### Added

- Make newError function public (#6). -@kostyantyn
- Add public access to default client (#5). -@kostyantyn
- Support default server mux in Handler.
- Allow error class to be customized from `honeybadger.Notify`.
- Support sending fingerprint in `honeybadger.Notify`.
- Added BeforeNotify callback.

### Fixed

- Drain the body of a response before closing it (#4). -@kostyantyn
- Update config at pointer rather than dereferencing. (#2).

## [0.0.1] - 2015-06-25

### Added

- Go client for Honeybadger.io.
