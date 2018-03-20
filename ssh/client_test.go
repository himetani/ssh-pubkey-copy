package ssh

import (
	"errors"
	"testing"
)

type MockSession struct {
	mockFunc func() error
	Session
}

func (m *MockSession) Exec(cmd string) ([]byte, error) { return nil, nil }

func (m *MockSession) Connect() ([]byte, error) {
	return nil, m.mockFunc()
}

func (m *MockSession) Close() {}

func TestPing(t *testing.T) {
	type data struct {
		TestName string
		MockFunc func() error
		Err      error
	}

	success := func() error { return nil }
	fail := func() error { return errors.New("failed") }

	tests := []data{
		{"Success", success, nil},
		{"Fail", fail, errors.New("failed")},
	}

	client := &PubKeyCopyClient{}

	for i, test := range tests {
		session := &MockSession{mockFunc: test.MockFunc}
		result := client.Ping(session)

		if test.Err == nil {
			if result != nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected nil, but got non-nil. msg=%s", i, test.TestName, result.Error())
			}
			continue
		}

		if test.Err != nil {
			if result == nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected non-nil, but got nil.", i, test.TestName)
				continue
			}

			if result.Error() != test.Err.Error() {
				t.Errorf("Test #%d %s (Error value Check): expected '%s', got '%s'", i, test.TestName, test.Err.Error(), result.Error())
			}
		}
	}
}
