package loggy

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var timestampRegexp = `[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(-[0-9]{2}:[0-9]{2}|Z)`

var logTestCases = []struct {
	Name                string
	Message             string
	Prefix              string
	Severity            int
	Threshold           int
	ExpectedStdoutRegex *regexp.Regexp
	ExpectedStderrRegex *regexp.Regexp
}{
	{
		Name:                "prefix",
		Message:             "what comes before?",
		Prefix:              "~~~",
		Severity:            LevelStd,
		Threshold:           LevelDebug,
		ExpectedStdoutRegex: regexp.MustCompile(timestampRegexp + ` OUT loggy.TestLogger_Log.func1 ~~~ what comes before\?`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-debug-threshold",
		Message:             "just testing",
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelDebug,
		ExpectedStdoutRegex: regexp.MustCompile(timestampRegexp + ` OUT loggy.TestLogger_Log.func1 just testing`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-standard-threshold",
		Message:             "just testing 2",
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: regexp.MustCompile(timestampRegexp + ` OUT loggy.TestLogger_Log.func1 just testing 2`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-critical-threshold",
		Message:             "just testing 3",
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelCritical,
		ExpectedStdoutRegex: regexp.MustCompile(timestampRegexp + ` OUT loggy.TestLogger_Log.func1 just testing 3`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "warning-severity-info-threshold",
		Message:             "something is not right",
		Prefix:              "",
		Severity:            LevelWarning,
		Threshold:           LevelInfo,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: regexp.MustCompile(timestampRegexp + ` WARN loggy.TestLogger_Log.func1 something is not right`),
	},
	{
		Name:                "critical-severity-critical-threshold",
		Message:             "BOOM",
		Prefix:              "",
		Severity:            LevelCritical,
		Threshold:           LevelCritical,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: regexp.MustCompile(timestampRegexp + ` CRIT loggy.TestLogger_Log.func1 BOOM`),
	},
	{
		Name:                "critical-severity-standard-threshold",
		Message:             "AAAAAAH ERRORRRR",
		Prefix:              "",
		Severity:            LevelCritical,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "warning-severity-standard-threshold",
		Message:             "Boo! Warning!",
		Prefix:              "",
		Severity:            LevelWarning,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "info-severity-standard-threshold",
		Message:             "Informative?",
		Prefix:              "",
		Severity:            LevelInfo,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "debug-severity-standard-threshold",
		Message:             "Informative?",
		Prefix:              "",
		Severity:            LevelDebug,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-negative-threshold",
		Message:             "Nothing to see here",
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           -1,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
}

func TestLogger_Log(t *testing.T) {
	for _, testCase := range logTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout := bytes.NewBuffer([]byte{})
			stderr := bytes.NewBuffer([]byte{})
			options := Options{
				Out:       stdout,
				Err:       stderr,
				Threshold: testCase.Threshold,
				Prefix:    testCase.Prefix,
			}
			l, ctx := New(context.Background(), options)
			err := l.Log(ctx, testCase.Severity, testCase.Message)
			assert.Nil(t, err)

			if (testCase.ExpectedStdoutRegex == nil && stdout.String() != "") ||
				(testCase.ExpectedStdoutRegex != nil && !testCase.ExpectedStdoutRegex.MatchString(stdout.String())) {

				t.Errorf("\ngot:            %q,\nexpected match: %s", stdout.String(), testCase.ExpectedStdoutRegex)
				return
			}

			if (testCase.ExpectedStderrRegex == nil && stderr.String() != "") ||
				(testCase.ExpectedStderrRegex != nil && !testCase.ExpectedStderrRegex.MatchString(stderr.String())) {

				t.Errorf("\ngot:            %q,\nexpected match: %s", stderr.String(), testCase.ExpectedStderrRegex)
				return
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	stdout := bytes.NewBuffer([]byte{})
	options := Options{
		Out:       stdout,
		Threshold: LevelInfo,
	}
	l, ctx := New(context.Background(), options)
	l.Info(ctx, "some message")

	regex := regexp.MustCompile(timestampRegexp + " INFO loggy.TestLogger_Info some message")
	assert.Regexp(t, regex, stdout.String())
}

func TestLogger_Tags(t *testing.T) {
	stdout := bytes.NewBuffer([]byte{})
	options := Options{
		Out:                 stdout,
		Threshold:           LevelInfo,
		DisableFunctionName: true,
		DisableTimestamps:   true,
	}
	l, ctx := New(context.Background(), options)

	var testTags = map[int]string{
		1: "waffles",
		2: "bacon",
		3: "waffles",
	}
	for i := 1; i <= len(testTags); i++ {
		_, ctx = l.AddTag(ctx, testTags[i], i)
	}

	tags := l.Tags(ctx)
	assert.Len(t, tags, 2)
	assert.Equal(t, 2, l.Tag(ctx, "bacon"))
	assert.Equal(t, 3, l.Tag(ctx, "waffles"))
}
