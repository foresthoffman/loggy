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
- Warning
  - Labeled "WARN".
  - Typically, indicates an issue that may require intervention.
- Info
  - Labeled "INFO".
  - Typically, indicates generic runtime information.
- Debug
  - Labeled "DEBUG".
  - Typically, indicates debug output.

For example, with a `Logger.Threshold` of `LevelCritical`, only logs the following severities would be output:

- Standard
  - Sent to stdout.
- Critical
  - Sent to stderr.

Likewise, if the `Logger.Threshold` was `LevelInfo`, all logs would be output except for those with a severity of Debug.

### Disabling Logging

Providing a `Logger.Threshold` < 0 will disable logging entirely. This behaves similarly to a standard `--quiet` CLI flag.

### Usage

```go
package main

import (
	"github.com/foresthoffman/loggy"
	"os"
	"bytes"
)

func main() {
	// - Use OS stdout/stderr
	// - Only show messages that are critical or standard.
	// - Custom prefix, prepended to each message before the timestamp.
	logger := loggy.New(os.Stdout, os.Stderr, "myPrefix", loggy.LevelCritical)
	
	// Send a standard message to stdout.
	logger.Std("hello!")
	
	// - Use custom stdout/stderr
	// - Only show messages that are information, warnings, critical, or standard.
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	logger = loggy.New(stdout, stderr, "", loggy.LevelInfo)
	
	// Send an error message, with the tag "error", to the custom stderr buffer.
	logger.Critical("something went wrong!", "error")

	// Send a debug message to the custom stdout buffer. This message will be ignored
	// because of the provided threshold.
	logger.Debug("get the fly swatter!")
}
```

### Testing

Run `go test -v -count=1 ./...` in the project root directory. Use the `-count=1` to force the tests to run un-cached.

_That's all, enjoy!_