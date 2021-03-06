package ssh

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type MockIn struct{}

func (m *MockIn) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (m *MockIn) Close() error {
	return nil
}

func TestPseudoTerminalSendTest(t *testing.T) {

	type data struct {
		TestName string
		In       io.WriteCloser
		Out      io.Reader
		Err      io.Reader
		Result   error
	}

	mockIn := &MockIn{}

	tests := []data{
		{"stdin is not initialized", mockIn, bytes.NewBuffer(nil), bytes.NewBuffer(nil), errors.New("Psedudo stdin is not initialized")},
	}

	for i, test := range tests {
		terminal := PseudoTerminal{}
		result := terminal.Send("")
		if test.Result == nil {
			if result != nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected nil, but got non-nil. msg=%s.", i, test.TestName, result.Error())
				continue
			}
		}

		if test.Result != nil {
			if result == nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected non-nil, but got nil.", i, test.TestName)
				continue
			}

			if result.Error() != test.Result.Error() {
				t.Errorf("Test #%d %s (Error value Check): expected '%s', got '%s'", i, test.TestName, test.Result.Error(), result.Error())
			}
		}
	}
}
