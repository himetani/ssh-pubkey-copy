package ssh

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	type data struct {
		TestName string
		Path     string
		Err      error
	}

	dummyPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "himetani", "ssh-pubkey-copy", "ssh", "testdata", "dummy_id_rsa")

	tests := []data{
		{"Invalid filepath", "hoge", errors.New("open hoge: no such file or directory")},
		{"Invalid file", dummyPath, errors.New("ssh: no key found")},
	}

	for i, test := range tests {
		_, err := NewPrivateKey(test.Path)
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

func TestNewPublicKeyContent(t *testing.T) {
	type data struct {
		TestName string
		Path     string
		Err      error
	}

	dummyPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "himetani", "ssh-pubkey-copy", "ssh", "testdata", "dummy_id_rsa.pub")

	tests := []data{
		{"Invalid filepath", "hoge", errors.New("open hoge: no such file or directory")},
		{"Valid file", dummyPath, nil},
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
