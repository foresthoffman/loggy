package loggy

import "testing"

var writeTestCases = []struct {
	Name            string
	Message         []byte
	WriteFn         WriteFn
	ExpectedMessage []byte
}{}

func TestWriter_Write(t *testing.T) {

}
