package ssh

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPing(t *testing.T) {
	privateKey := filepath.Join(os.Getenv("GOPATH"), "src/github.com/himetani/ssh-pubkey-copy/.vagrant/machines/default/virtualbox/private_key")

	type data struct {
		TestName  string
		Dests     []Dest
		ResultErr error
	}

	tests := []data{
		{"Copied", []Dest{Dest{Host: "127.0.0.1", Port: "2222", User: "vagrant"}}, nil},
		{"Not Copied", []Dest{Dest{Host: "127.0.0.1", Port: "2222", User: "clienttest1"}}, errors.New("ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain")},
	}

	client := KeyClient{
		privateKey: privateKey,
	}

	for i, test := range tests {
		results := client.Ping(test.Dests)
		if len(results) != 1 {
			t.Errorf("Test #%d %s: Unexpected error happened. len = %d\n", i, test.TestName, len(results))
		}

		result := results[0]
		if test.ResultErr == nil {
			if result.Err != nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected nil, but got non-nil. msg=%s", i, test.TestName, result.Err.Error())
			}
			continue
		}

		if test.ResultErr != nil {
			if result.Err == nil {
				t.Errorf("Test #%d %s: Unexpected error happend. err expected non-nil, but got nil.", i, test.TestName)
				continue
			}

			if result.Err.Error() != test.ResultErr.Error() {
				t.Errorf("Test #%d %s (ResultError Check): expected '%s', got '%s'", i, test.TestName, test.ResultErr.Error(), result.Err.Error())
			}
		}

	}
}
