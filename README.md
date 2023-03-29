# 🌲 Loggy 🌲

![Current project build status](https://github.com/foresthoffman/loggy/actions/workflows/go.yml/badge.svg)

Loggy wraps the standard log.Logger structure with some additional fields, to standardize stdout and stderr output. Based on the "severity" of any given log, loggy redirects logs to stdout while redirecting errors to stderr. The desired logging level, or "threshold", shared between the two output streams filters the output messages. If a log or error message doesn't fall on or under the desired threshold, then that message is ignored. Messages with the "standard" severity are always displayed.

### Installation

Run `go get -u github.com/foresthoffman/loggy`

If you're using `go mod`, run `go mod vendor` afterwards.

### Importing

Import this package by including `github.com/foresthoffman/loggy` in your import block.

e.g.

```go
package main

import(
    ...
    "github.com/foresthoffman/loggy"
)
```

### Log Levels

The logger's `Threshold` is an integer that determines what severities (e.g. "log levels") are output and which are ignored. Users may assign whatever context and importance they deem fit to the log levels; however, Critical-level and Warning-level logs will be output to the provided stderr regardless. In descending order the available log levels are:

- Standard
  - Labeled "OUT".
  - Typically, just standard output of no particular importance.
- Critical
  - Labeled "CRIT".
  - Typically, indicates a fatal issue.
- Error
  - Labeled "ERROR".
  - Typically, indicates a general issue that was recovered.
- Warning
  - Labeled "WARN".
  - Typically, indicates an issue that may require intervention.
- Info
  - Labeled "INFO".
  - Typically, indicates generic runtime information.
- Debug
  - Labeled "DEBUG".
  - Typically, indicates debug output.

For example, with a threshold of `LevelCritical`, only logs the following severities would be output:

- Standard
  - Sent to stdout.
- Critical
  - Sent to stderr.

Likewise, if the threshold was `LevelInfo`, all logs would be output except for those with a severity of Debug.

### Disabling Logging

Providing a `Logger.Threshold` < 0 will disable logging entirely. This behaves similarly to a standard `--quiet` CLI flag.

### Usage

```go
package main

import (
  "bytes"
  "context"
  "errors"
  "fmt"
  "github.com/foresthoffman/loggy"
)

func main() {
  // - Use OS stdout/stderr by default.
  // - Only show messages that are critical or standard.
  // - Custom prefix, prepended to each message before the timestamp.
  options := loggy.Options{
    Prefix: "~~~",
    Threshold: loggy.LevelCritical,
  }
  logger, ctx := loggy.New(context.Background(), options)

  // Send a standard message to stdout.
  logger.Std(ctx, "hello!") // 2023-03-29T15:20:55-05:00 OUT main.main ~~~ hello!

  // - Use custom stdout/stderr
  // - Only show messages that are information, warnings, critical, or standard.
  stdout := bytes.NewBuffer([]byte{})
  stderr := bytes.NewBuffer([]byte{})
  options = loggy.Options{
    Out: stdout,
    Err: stderr,
    Threshold: loggy.LevelInfo,
  }
  logger, ctx = loggy.New(context.Background(), options)

  // Send an error message, with the tag "error", to the custom stderr buffer.
  err := errors.New("oops")
  logger.Critical(ctx, "something went wrong!", err.Error())

  // Send a debug message to the custom stdout buffer. This message will be ignored
  // because of the provided threshold.
  logger.Debug(ctx, "get the fly swatter!")

  fmt.Println(stderr) // 2023-03-29T15:25:27-05:00 CRIT main.main something went wrong! oops
}
```

### Testing

Run `go test -v -count=1 ./...` in the project root directory. Use the `-count=1` to force the tests to run un-cached.

_That's all, enjoy!_