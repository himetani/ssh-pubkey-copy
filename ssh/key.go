package ssh

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

func NewPrivateKey(path string) (ssh.Signer, error) {
	if path == "" {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".ssh", "id_rsa")
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(buf)
}

func NewPublicKeyContent(path string) (string, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(buf), "\n"), nil
}
