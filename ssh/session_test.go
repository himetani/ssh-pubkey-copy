package ssh

import (
	"testing"
)

const (
	ValidIP     = "localhost"
	ValidPort   = "2222"
	ValidUser   = "vagrant"
	ValidPasswd = "vagrant"
)

func TestNewPseudoTerminal(t *testing.T) {
	type data struct {
		TestName string
		IP       string
		Port     string
		User     string
		Passwd   string
		Err      error
	}
	tests := []data{
		{"Valid params", ValidIP, ValidPort, ValidUser, ValidPasswd, nil},
	}

	for i, test := range tests {
		_, err := NewPseudoTerminal(test.IP, test.Port, test.User, test.Passwd)
		if test.Err == nil {
			if err != nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected nil, but got non-nil. msg=%s.", i, test.TestName, err.Error())
				continue
			}
		}

		if test.Err != nil {
			if err == nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected non-nil, but got nil.", i, test.TestName)
				continue
			}

			if err.Error() != test.Err.Error() {
				t.Errorf("Test #%d %s (Error value Check): expected '%s', got '%s'", i, test.TestName, test.Err.Error(), err.Error())
			}
		}
	}
}
