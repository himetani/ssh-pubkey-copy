package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Session is the interface that wraps session behavior.
type Session interface {
	Executor
	Connector
	Closer
}

// Executor is the interface that wraps Exec method
type Executor interface {
	Exec(string) ([]byte, error)
}

// Connector is the interface that wraps Connect method
type Connector interface {
	Connect() ([]byte, error)
}

// Closer is the interface that wraps Close method
type Closer interface {
	Close()
}

// Sender is the interface that wraps Send method
type Sender interface {
	Send(string) error
}

// Starter is the interface that wraps Start method
type Starter interface {
	Start() error
}

// UserSwitcher is the interface that wraps SwitchUser method
type UserSwitcher interface {
	SwitchUser(string, string) error
}

// Ender is the interface that wraps Start method
type Ender interface {
	End() error
}

// NewPrivateKeySession returns Session Wrapper instance whose session is established by private key
func NewPrivateKeySession(ip, port, user string, privateKey ssh.Signer) (*Wrapper, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
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

	return &Wrapper{
		conn:    conn,
		session: session,
	}, nil
}

// NewPasswordSession returns Session Wrapper instance whose session is established by password
func NewPasswordSession(ip, port, user, password string) (*Wrapper, error) {
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

	return &Wrapper{
		conn:    conn,
		session: session,
	}, nil
}

// Wrapper is the struct that wraps ssh connectivity
type Wrapper struct {
	conn    *ssh.Client
	session *ssh.Session
	Session
}

// Close is the function to close the session & connection
func (s *Wrapper) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Connect is the funcion to connect using seession
func (w *Wrapper) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return w.session.Output(cmd)
}

// Exec is the function to exec any cmd on the session
func (w *Wrapper) Exec(cmd string) ([]byte, error) {
	return w.session.Output(cmd)
}
