package loggy

import (
	"io"
)

var _ io.Writer = &Writer{}

// WriteFn is a handler for processing each individual byte slice before it gets
// sent to the target stream.
type WriteFn = func(out io.Writer, p []byte) error

type Writer struct {
	handler WriteFn
	out     io.Writer
}

func DefaultWriteFn(out io.Writer, p []byte) error {
	if _, err := out.Write(p); err != nil {
		return err
	}
	return nil
}

func NewWriter(out io.Writer, fn WriteFn) *Writer {
	return &Writer{
		handler: fn,
		out:     out,
	}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	err = w.handler(w.out, p)
	if err != nil {
		return
	}

	n = len(p)
	return
}
