package loggy

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

var writeTestCases = []struct {
	Name            string
	Message         []byte
	WriteFn         WriteFn
	ExpectedMessage []byte
}{
	{
		Name:            "basic",
		Message:         []byte("Hello, basic!"),
		WriteFn:         DefaultWriteFn,
		ExpectedMessage: []byte("Hello, basic!"),
	},
	{
		Name:            "empty",
		Message:         []byte(""),
		WriteFn:         DefaultWriteFn,
		ExpectedMessage: []byte(""),
	},
	{
		Name:    "write-fn-info",
		Message: []byte("ERROR something went wrong\nINFO interesting!"),
		WriteFn: func(out io.Writer, p []byte) error {
			lines := bytes.Split(p, []byte("\n"))
			for _, line := range lines {
				if !bytes.Contains(line, []byte("INFO")) {
					continue
				}
				n, err := out.Write(line)
				if err != nil {
					return err
				}
				if n < len(line) {
					return errors.New("must write whole line")
				}
			}
			return nil
		},
		ExpectedMessage: []byte("INFO interesting!"),
	},
	{
		Name:    "write-fn-error",
		Message: []byte("ERROR something went wrong\nINFO interesting!"),
		WriteFn: func(out io.Writer, p []byte) error {
			lines := bytes.Split(p, []byte("\n"))
			for _, line := range lines {
				if !bytes.Contains(line, []byte("ERROR")) {
					continue
				}
				n, err := out.Write(line)
				if err != nil {
					return err
				}
				if n < len(line) {
					return errors.New("must write whole line")
				}
			}
			return nil
		},
		ExpectedMessage: []byte("ERROR something went wrong"),
	},
}

func TestWriter_Write(t *testing.T) {
	for _, testCase := range writeTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			stdout := bytes.NewBuffer([]byte{})
			w := NewWriter(stdout, testCase.WriteFn)

			n, err := w.Write(testCase.Message)
			if err != nil {
				t.Error(err)
				return
			}
			if n < len(testCase.Message) {
				t.Errorf("got: %d, expected: %d", n, len(testCase.Message))
				return
			}

			if stdout.String() != string(testCase.ExpectedMessage) {
				t.Errorf("\ngot:      %q,\nexpected: %q", stdout.String(), string(testCase.ExpectedMessage))
				return
			}
		})
	}
}
