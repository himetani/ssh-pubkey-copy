package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Terminal is the interface that wraps session behavior and terminal specific behavior.
type Terminal interface {
	Session
	Starter
	Sender
	UserSwitcher
	Ender
}

// PseudoTerminal is the struct that represents psudo terminal ssh session
type PseudoTerminal struct {
	conn    *ssh.Client
	session *ssh.Session
	Terminal
	in  io.WriteCloser
	out io.Reader
	err io.Reader
}

// Close is the function to close the session & connection
func (s *PseudoTerminal) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Connect is the funcion to connect using seession
func (w *PseudoTerminal) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return w.session.Output(cmd)
}

// Exec is the function to exec any cmd on the session
func (p *PseudoTerminal) Exec(cmd string) ([]byte, error) {
	return p.session.Output(cmd)
}

// Start is the function to start session as a pseudo terminal
func (p *PseudoTerminal) Start() error {
	defer fmt.Println("[Debug] Start Pseudo Terminal")

	out := bytes.NewBuffer(nil)
	err := bytes.NewBuffer(nil)

	p.in, _ = p.session.StdinPipe()
	p.out, _ = p.session.StdoutPipe()
	p.err, _ = p.session.StderrPipe()

	p.session.Stdout = out
	p.session.Stderr = err

	if err := p.session.Shell(); err != nil {
		return errors.New(fmt.Sprintf("failed to start shell: %s", err))
	}

	return nil
}

// Send is the function to send cmd to a pseudo terminal
func (p *PseudoTerminal) Send(cmd string) error {
	defer fmt.Println("[Debug] Send cmd to Pseudo Terminal")

	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.session.Stdout == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.session.Stderr == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	cmd = fmt.Sprintf("%s; echo end\n", cmd)

	fmt.Fprintf(p.in, cmd)

	scanner := bufio.NewScanner(p.out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "end") {
			time.Sleep(1 * time.Second)
			break
		}
	}

	return nil
}

// SwitchUser is the function to switch user on pseudo terminal login sesion
func (p *PseudoTerminal) SwitchUser(user, passwd string) error {
	defer fmt.Printf("[Debug] Switch User to Pseudo Terminal. Username: %s\n", user)

	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.session.Stdout == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.session.Stderr == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		scanner := bufio.NewScanner(p.out)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "password") {
				fmt.Printf("[Debug] End of scanning password line: %s\n", line)
				wg.Done()
				return
			}
		}
	}()

	fmt.Println("[Debug] Start to input command")
	sudoCmd := fmt.Sprintf("sudo su - %s\n", user)
	fmt.Fprintf(p.in, sudoCmd)

	time.Sleep(100 * time.Millisecond)
	inputPasswd := fmt.Sprintf("%s\n", passwd)
	fmt.Fprintf(p.in, inputPasswd)

	wg.Wait()
	return nil
}

// End is the function to end session as a pseudo terminal
func (p *PseudoTerminal) End() error {
	defer fmt.Println("[Debug] End Pseudo Terminal")

	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.session.Stdout == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.session.Stderr == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	fmt.Fprintf(p.in, "%s\n", "exit")
	return nil
}

// NewPseudoTerminal returns PseudoTerminal instance whose session is established by password.
func NewPseudoTerminal(ip, port, user, password string) (*PseudoTerminal, error) {
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

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, errors.New(fmt.Sprintf("request for pseudo terminal failed: %s", err))
	}

	return &PseudoTerminal{
		conn:    conn,
		session: session,
	}, nil
}
