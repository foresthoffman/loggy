package loggy

import (
	"bytes"
	"regexp"
	"testing"
)

var logTestCases = []struct {
	Name                string
	Message             string
	Tags                []string
	Prefix              string
	Severity            int
	Threshold           int
	ExpectedStdoutRegex *regexp.Regexp
	ExpectedStderrRegex *regexp.Regexp
}{
	{
		Name:                "prefix",
		Message:             "what comes before?",
		Tags:                nil,
		Prefix:              "~~~",
		Severity:            LevelStd,
		Threshold:           LevelDebug,
		ExpectedStdoutRegex: regexp.MustCompile(`~~~[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} OUT testing.tRunner \[] what comes before\?\n`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-debug-threshold",
		Message:             "just testing",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelDebug,
		ExpectedStdoutRegex: regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} OUT testing.tRunner \[] just testing\n`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-standard-threshold",
		Message:             "just testing 2",
		Tags:                []string{"super standard"},
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} OUT testing.tRunner \[super standard] just testing 2\n`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "standard-severity-critical-threshold",
		Message:             "just testing 3",
		Tags:                []string{"still", "very", "standard"},
		Prefix:              "",
		Severity:            LevelStd,
		Threshold:           LevelCritical,
		ExpectedStdoutRegex: regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} OUT testing.tRunner \[still, very, standard] just testing 3\n`),
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "warning-severity-info-threshold",
		Message:             "something is not right",
		Tags:                []string{"warning"},
		Prefix:              "",
		Severity:            LevelWarning,
		Threshold:           LevelInfo,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} WARN testing.tRunner \[warning] something is not right\n`),
	},
	{
		Name:                "critical-severity-critical-threshold",
		Message:             "BOOM",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelCritical,
		Threshold:           LevelCritical,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: regexp.MustCompile(`[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} CRIT testing.tRunner \[] BOOM\n`),
	},
	{
		Name:                "critical-severity-standard-threshold",
		Message:             "AAAAAAH ERRORRRR",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelCritical,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "warning-severity-standard-threshold",
		Message:             "Boo! Warning!",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelWarning,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "info-severity-standard-threshold",
		Message:             "Informative?",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelInfo,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
	{
		Name:                "debug-severity-standard-threshold",
		Message:             "Informative?",
		Tags:                nil,
		Prefix:              "",
		Severity:            LevelDebug,
		Threshold:           LevelStd,
		ExpectedStdoutRegex: nil,
		ExpectedStderrRegex: nil,
	},
}

func TestLogger_Log(t *testing.T) {
	for _, testCase := range logTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout := bytes.NewBuffer([]byte{})
			stderr := bytes.NewBuffer([]byte{})
			logger := New(stdout, stderr, testCase.Prefix, testCase.Threshold)
			logger.Log(testCase.Severity, testCase.Message, testCase.Tags...)

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
