package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

// Session is struct representing ssh Session
type Session struct {
	config    *ssh.ClientConfig
	conn      *ssh.Client
	session   *ssh.Session
	StdinPipe io.WriteCloser
}

// NewPrivateKeySession returns new Session instance
func NewPrivateKeySession(ip, port, user, privateKey string) (*Session, error) {
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

	return &Session{
		config:  config,
		conn:    conn,
		session: session,
	}, nil
}

// NewPasswordSession returns new Session instance
func NewPasswordSession(ip, port, user, password string) (*Session, error) {
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

	return &Session{
		config:  config,
		conn:    conn,
		session: session,
	}, nil
}

// Close close the session & connection
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Connect is func to connect
func (s *Session) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return s.session.Output(cmd)
}

// Exec is func to exec cmd on the session
func (s *Session) Exec(cmd string) ([]byte, error) {
	return s.session.Output(cmd)
}
