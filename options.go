package loggy

import (
	"io"
	"os"
	"time"
)

type Options struct {
	// The underlying stdout logger.
	Out io.Writer
	// The underlying stderr logger.
	Err io.Writer
	// The maximum severity to display for this logger. To disable logging completely, provide a Level < 0.
	Threshold Level
	// The text to place at the beginning of each log message, after the timestamp,
	// severity, function name, and context tags.
	Prefix string
	// Set to true to disable timestamps. This is useful if piping logs into a writer
	// that already uses timestamps.
	DisableTimestamps bool
	// Time format to use to output timestamps.
	TimestampFormat string
	// Timestamp function to get current time.
	TimestampFunc func() time.Time
	// Set to true to log un-resolvable internal errors as fatal logs. Otherwise, return the errors and log nothing.
	LogFatal bool
	// Set to true to disable outputting the calling function name before the rest of the log message.
	DisableFunctionName bool
	// Set to true to disable outputting the context tags. This purely hides the tag
	// list from being prepended to any log messages, the *Tag* helper functions will
	// still work and will still manage state.
	DisableTags bool
	// The context key where the logger can store tags exposed by the *Tag* helper functions.
	TagsContextKey string
}

// DefaultOptions contains all the standard options that a logger will use when certain options are not provided.
var DefaultOptions = Options{
	Out:                 os.Stdout,
	Err:                 os.Stderr,
	Threshold:           LevelInfo,
	Prefix:              "",
	DisableTimestamps:   false,
	TimestampFormat:     time.RFC3339,
	TimestampFunc:       time.Now,
	LogFatal:            false,
	DisableFunctionName: false,
	TagsContextKey:      ContextKeyTags,
}
