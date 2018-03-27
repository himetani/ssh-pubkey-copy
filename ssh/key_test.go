package ssh

import (
	"errors"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	type data struct {
		TestName string
		Path     string
		Err      error
	}

	tests := []data{
		{"Invalid file", "hoge", errors.New("open hoge: no such file or directory")},
	}

	for i, test := range tests {
		_, err := NewPrivateKey(test.Path)
		if test.Err == nil {
			continue
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

func TestNewPublicKeyContent(t *testing.T) {
	type data struct {
		TestName string
		Path     string
		Err      error
	}

	tests := []data{
		{"Invalid file", "hoge", errors.New("open hoge: no such file or directory")},
	}

	for i, test := range tests {
		_, err := NewPrivateKey(test.Path)
		if test.Err == nil {
			continue
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
