// Package loggy wraps the standard log.Logger with varying levels of verbosity
// and the option to split stdout/stderr streams for info/error logs. The default
// level of verbosity is always shown, to prevent any invalid logs from being
// lost by accident.
package loggy

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
)

const (
	// ContextKeyLogger is the context.Context key where loggy logger references are stored.
	ContextKeyLogger = "loggy.Logger"
	// ContextKeyTags is the context.Context key where loggy tags are stored.
	ContextKeyTags = "loggy.Tags"
)

// Must implement interface.
var _ Logger = &logger{}

type Logger interface {
	Log(ctx context.Context, severity Level, message ...interface{}) error
	Logf(ctx context.Context, severity Level, format string, message ...interface{}) error
	Std(ctx context.Context, message ...interface{}) error
	Stdf(ctx context.Context, format string, message ...interface{}) error
	Critical(ctx context.Context, message ...interface{}) error
	Criticalf(ctx context.Context, format string, message ...interface{}) error
	Warning(ctx context.Context, message ...interface{}) error
	Warningf(ctx context.Context, format string, message ...interface{}) error
	Info(ctx context.Context, message ...interface{}) error
	Infof(ctx context.Context, format string, message ...interface{}) error
	Debug(ctx context.Context, message ...interface{}) error
	Debugf(ctx context.Context, format string, message ...interface{}) error
	Tags(ctx context.Context) map[string]interface{}
	Tag(ctx context.Context, name string) interface{}
	AddTag(ctx context.Context, name string, value interface{}) (map[string]interface{}, context.Context)
	RemoveTag(ctx context.Context, name string) (map[string]interface{}, context.Context)
}

type logger struct {
	options *Options
	mux     sync.Mutex

	Ctx context.Context
}

// New creates a new wrapper for the log.Logger standard package. The provided
// threshold determines what level of verbosity the provided stream will receive.
func New(ctx context.Context, options Options) (*logger, context.Context) {
	l := &logger{
		options: &options,
	}
	if l.options.Out == nil {
		l.options.Out = DefaultOptions.Out
	}
	if l.options.Err == nil {
		l.options.Err = DefaultOptions.Err
	}
	if l.options.TimestampFormat == "" {
		l.options.TimestampFormat = DefaultOptions.TimestampFormat
	}
	if l.options.TimestampFunc == nil {
		l.options.TimestampFunc = DefaultOptions.TimestampFunc
	}
	if l.options.TagsContextKey == "" {
		l.options.TagsContextKey = DefaultOptions.TagsContextKey
	}

	return l, context.WithValue(ctx, ContextKeyLogger, l)
}

// Log is a wrapper for Logf without the format string.
func (l *logger) Log(ctx context.Context, severity Level, message ...interface{}) error {
	return l.Logf(ctx, severity, "", message...)
}

// Logf gathers the provided message metadata and writes the compiled message to
// the configured output or error stream, depending on severity. By default, log
// messages are prefixed with: a timestamp, log severity, log function name, and
// any tags assigned to the context via the *Tag* helper methods. All of these
// features can be figured via loggy.Options, when using loggy.New().
func (l *logger) Logf(ctx context.Context, severity Level, format string, message ...interface{}) error {
	if l.options.Threshold < 0 {
		// Logging is disabled.
		return nil
	}
	if severity < 0 || severity+1 > len(LevelNames) {
		severity = LevelStd
	}
	if severity != LevelStd && severity > l.options.Threshold {
		return nil
	}
	var msg = fmt.Sprintf("%s", LevelNames[severity])

	if !l.options.DisableFunctionName {
		// Get calling function name.
		pc, _, _, ok := runtime.Caller(2)
		if !ok {
			lookupErr := fmt.Sprintf(
				"%s %s %s",
				LevelNames[LevelCritical],
				"loggy.logger.Logf",
				"failed to dynamically lookup function name",
			)
			_, err := l.options.Err.Write([]byte(l.maybePrefixTimestamp(lookupErr)))
			if err != nil {
				if l.options.LogFatal {
					log.Fatal(lookupErr)
				} else {
					return err
				}
			}
		} else {
			fullName := strings.Split(runtime.FuncForPC(pc).Name(), "/")

			msg = fmt.Sprintf("%s %s", msg, fullName[len(fullName)-1])
		}
	}

	if !l.options.DisableTags {
		// Compile tags from context.
		tags := l.Tags(ctx)
		count := 0
		if tags != nil && len(tags) > 0 {
			tagBytes := []byte("[")
			for name, value := range tags {
				var delim string
				if count+1 < len(tags) {
					delim = ", "
				}
				tagBytes = append(tagBytes, []byte(fmt.Sprintf("%s:%v%s", name, value, delim))...)
				count++
			}
			tagBytes = append(tagBytes, []byte("]")...)

			msg = fmt.Sprintf("%s %s", msg, string(tagBytes))
		}
	}

	if l.options.Prefix != "" {
		// Append prefix before the user-formatted message.
		msg = fmt.Sprintf("%s %s", msg, l.options.Prefix)
	}

	// Append user-formatted message.
	if format == "" && len(message) > 0 {
		for i := 0; i < len(message); i++ {
			format = format + " %v"
		}
	}
	message = append([]interface{}{msg}, message...)
	msg = fmt.Sprintf("%s"+format+"\n", message...)
	if severity == LevelStd || severity >= LevelInfo {
		_, err := l.options.Out.Write([]byte(l.maybePrefixTimestamp(msg)))
		if err != nil {
			if l.options.LogFatal {
				log.Fatal(msg)
			} else {
				return err
			}
		}
	} else {
		_, err := l.options.Err.Write([]byte(l.maybePrefixTimestamp(msg)))
		if err != nil {
			if l.options.LogFatal {
				log.Fatal(msg)
			} else {
				return err
			}
		}
	}

	return nil
}

// Std sends a standard log message.
func (l *logger) Std(ctx context.Context, message ...interface{}) error {
	return l.Logf(ctx, LevelStd, "", message...)
}

// Stdf sends a standard log message, with a custom string format.
func (l *logger) Stdf(ctx context.Context, format string, message ...interface{}) error {
	return l.Logf(ctx, LevelStd, format, message...)
}

// Critical sends a critical error message.
func (l *logger) Critical(ctx context.Context, message ...interface{}) error {
	return l.Logf(ctx, LevelCritical, "", message...)
}

// Criticalf sends a critical error message, with a custom string format.
func (l *logger) Criticalf(ctx context.Context, format string, message ...interface{}) error {
	return l.Logf(ctx, LevelCritical, format, message...)
}

// Warning sends a warning error message.
func (l *logger) Warning(ctx context.Context, message ...interface{}) error {
	return l.Logf(ctx, LevelWarning, "", message...)
}

// Warningf sends a warning error message, with a custom string format.
func (l *logger) Warningf(ctx context.Context, format string, message ...interface{}) error {
	return l.Logf(ctx, LevelWarning, format, message...)
}

// Info sends an info log message.
func (l *logger) Info(ctx context.Context, message ...interface{}) error {
	return l.Logf(ctx, LevelInfo, "", message...)
}

// Infof sends an info log message, with a custom string format.
func (l *logger) Infof(ctx context.Context, format string, message ...interface{}) error {
	return l.Logf(ctx, LevelInfo, format, message...)
}

// Debug sends a debug log message.
func (l *logger) Debug(ctx context.Context, message ...interface{}) error {
	return l.Logf(ctx, LevelDebug, "", message...)
}

// Debugf sends a debug log message, with a custom string format.
func (l *logger) Debugf(ctx context.Context, format string, message ...interface{}) error {
	return l.Logf(ctx, LevelDebug, format, message...)
}

// Tags returns all tags associated with the provided context.
func (l *logger) Tags(ctx context.Context) map[string]interface{} {
	l.mux.Lock()
	defer l.mux.Unlock()

	tags, ok := ctx.Value(l.options.TagsContextKey).(map[string]interface{})
	if !ok {
		tags = make(map[string]interface{})
	}

	return tags
}

// Tag returns an individual tag, by name, associated with the provided context.
func (l *logger) Tag(ctx context.Context, name string) interface{} {
	l.mux.Lock()
	defer l.mux.Unlock()

	tags, ok := ctx.Value(l.options.TagsContextKey).(map[string]interface{})
	if !ok {
		return nil
	}
	tag, ok := tags[name]
	if !ok {
		return nil
	}

	return tag
}

// AddTag adds or updates a tag, by name, associated with the provided context.
//
// NOTE: Be wary of adding tags in any goroutines if there's any possibility of
// duplicate tag names. Although loggy uses mutexes to ensure there's no race
// condition, there's no guarantee as to which value will be stored last.
func (l *logger) AddTag(ctx context.Context, name string, value interface{}) (map[string]interface{}, context.Context) {
	l.mux.Lock()
	defer l.mux.Unlock()

	tags, ok := ctx.Value(l.options.TagsContextKey).(map[string]interface{})
	if !ok {
		tags = make(map[string]interface{})
	}
	if name != "" {
		tags[name] = value
		ctx = context.WithValue(ctx, l.options.TagsContextKey, tags)
	}

	return tags, ctx
}

// RemoveTag removes a tag, by name, associated with the provided context.
func (l *logger) RemoveTag(ctx context.Context, name string) (map[string]interface{}, context.Context) {
	l.mux.Lock()
	defer l.mux.Unlock()

	tags, ok := ctx.Value(l.options.TagsContextKey).(map[string]interface{})
	if !ok {
		tags = make(map[string]interface{})
	}
	if name != "" {
		delete(tags, name)
		ctx = context.WithValue(ctx, l.options.TagsContextKey, tags)
	}

	return tags, ctx
}

func (l *logger) maybePrefixTimestamp(msg string) string {
	if !l.options.DisableTimestamps {
		msg = fmt.Sprintf(
			"%s %s",
			l.options.TimestampFunc().
				Format(l.options.TimestampFormat), msg)
	}
	return msg
}
