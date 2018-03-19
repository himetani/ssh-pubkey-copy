package ssh

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHSession is struct representing ssh Session
type SSHSession struct {
	conn    *ssh.Client
	session *ssh.Session
}

type Session interface {
	Executor
	Connector
	Closer
}

type Executor interface {
	Exec()
}

type Connector interface {
	Connect()
}

type Closer interface {
	Close()
}

// NewPrivateKeySession returns new Session instance
func NewPrivateKeySession(ip, port, user, privateKey string) (*SSHSession, error) {
	buf, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: 1 * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return &SSHSession{
		conn:    conn,
		session: session,
	}, nil
}

// NewPasswordSession returns new Session instance
func NewPasswordSession(ip, port, user, password string) (*SSHSession, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 1 * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return &SSHSession{
		conn:    conn,
		session: session,
	}, nil
}

// Close close the session & connection
func (s *SSHSession) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Connect is func to connect
func (s *SSHSession) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return s.session.Output(cmd)
}

// Exec is func to exec cmd on the session
func (s *SSHSession) Exec(cmd string) ([]byte, error) {
	return s.session.Output(cmd)
}

type Result struct {
	*Dest
	Err error
}
