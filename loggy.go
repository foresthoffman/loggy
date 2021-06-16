// Package loggy wraps the standard log.Logger with varying levels of verbosity
// and the option to split stdout/stderr streams for info/error logs. The default
// level of verbosity is always shown, to prevent any invalid logs from being
// lost by accident.
//
// ## Split stdout and stderr:
//
// logger := loggy.New(os.Stdout, os.Stderr, "myPrefix", loggy.LevelInfo)
//
// ## Combined stdout and stderr:
//
// logger := loggy.NewCombined(os.Stdout, "myPrefix", loggy.LevelInfo)
package loggy

import (
	"io"
	"log"
	"runtime"
	"strings"
)

type Logger struct {
	// The underlying stdout logger.
	Stdout *log.Logger
	// The underlying stderr logger.
	Stderr *log.Logger
	// The maximum severity to display for this logger.
	Threshold int
}

const (
	// LevelStd indicates standard log output. Always shown.
	LevelStd = iota
	// LevelCritical indicates a fatal issue.
	LevelCritical
	// LevelWarning indicates an issue that may require intervention.
	LevelWarning
	// LevelInfo indicates generic runtime information.
	LevelInfo
	// LevelDebug indicates debug output.
	LevelDebug
)

// LevelNames describe the alphabetical types to label each Level* with in stdout/stderr.
var LevelNames = []string{"OUT", "CRIT", "WARN", "INFO", "DEBUG"}

// New creates a new wrapper for the log.Logger standard package. The provided
// threshold determines what level of verbosity the provided stream will receive.
// Default threshold only captures standard output and critical errors.
func New(stdout, stderr io.Writer, prefix string, threshold int) *Logger {
	return &Logger{
		Stdout:    log.New(stdout, prefix, log.LstdFlags),
		Stderr:    log.New(stderr, prefix, log.LstdFlags),
		Threshold: threshold,
	}
}

func NewCombined(out io.Writer, prefix string, threshold int) *Logger {
	return New(out, out, prefix, threshold)
}

// Log gathers the provided message metadata and tries to capture the name of the
// calling function, dynamically. It then writes the compiled message to either
// standard out or standard error, depending on severity. If the name of the
// calling function cannot be captured, and the severity happens to be set to
// debug, an internal debug log will be sent accordingly.
func (l *Logger) Log(severity int, message string, tags ...string) {
	if l.Threshold < 0 {
		return
	}
	if severity < 0 || severity+1 > len(LevelNames) {
		severity = LevelStd
	}
	if severity != LevelStd && severity > l.Threshold {
		return
	}

	var fn string
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		if severity == LevelDebug {
			l.Logf(LevelDebug, "loggy.Logger.Log", "failed to dynamically lookup function name")
		}
	} else {
		fullName := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		fn = fullName[len(fullName)-1]
	}

	l.Logf(severity, fn, message, tags...)
}

// Logf gathers the provided message metadata and writes the compiled message to
// either standard out or standard error, depending on severity.
func (l *Logger) Logf(severity int, fn, message string, tags ...string) {
	if l.Threshold < 0 {
		return
	}
	if severity < 0 || severity+1 > len(LevelNames) {
		severity = LevelStd
	}
	if severity != LevelStd && severity > l.Threshold {
		return
	}

	tagsList := []byte("[")
	for i, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		tagsList = append(tagsList, []byte(tag)...)
		if i+1 < len(tags) {
			tagsList = append(tagsList, []byte(", ")...)
		}
	}
	tagsList = append(tagsList, ']')

	if severity == LevelStd || severity >= LevelInfo {
		l.Stdout.Println(LevelNames[severity], fn, string(tagsList), message)
	} else {
		l.Stderr.Println(LevelNames[severity], fn, string(tagsList), message)
	}
}

// Std sends a standard log message.
func (l *Logger) Std(message string, tags ...string) {
	l.Log(LevelStd, message, tags...)
}

// Stdf sends a standard log message, with a custom function name.
func (l *Logger) Stdf(fn, message string, tags ...string) {
	l.Logf(LevelStd, fn, message, tags...)
}

// Critical sends a critical error message.
func (l *Logger) Critical(message string, tags ...string) {
	l.Log(LevelCritical, message, tags...)
}

// Criticalf sends a critical error message, with a custom function name.
func (l *Logger) Criticalf(fn, message string, tags ...string) {
	l.Logf(LevelCritical, fn, message, tags...)
}

// Warning sends a warning error message.
func (l *Logger) Warning(message string, tags ...string) {
	l.Log(LevelWarning, message, tags...)
}

// Warningf sends a warning error message, with a custom function name.
func (l *Logger) Warningf(fn, message string, tags ...string) {
	l.Logf(LevelWarning, fn, message, tags...)
}

// Info sends a info log message.
func (l *Logger) Info(message string, tags ...string) {
	l.Log(LevelInfo, message, tags...)
}

// Infof sends a info log message, with a custom function name.
func (l *Logger) Infof(fn, message string, tags ...string) {
	l.Logf(LevelInfo, fn, message, tags...)
}

// Debug sends a debug log message.
func (l *Logger) Debug(message string, tags ...string) {
	l.Log(LevelDebug, message, tags...)
}

// Debugf sends a debug log message, with a custom function name.
func (l *Logger) Debugf(fn, message string, tags ...string) {
	l.Logf(LevelDebug, fn, message, tags...)
}
