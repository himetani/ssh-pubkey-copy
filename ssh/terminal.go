package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// PseudoTerminal is the struct that represents psudo terminal ssh session
type PseudoTerminal struct {
	wrapper *Wrapper
	in      io.WriteCloser
	out     io.Reader
	err     io.Reader
}

// Close is the function to close the session & connection
func (p *PseudoTerminal) Close() {
	p.wrapper.Close()
}

// Connect is the funcion to connect using seession
func (p *PseudoTerminal) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return p.Exec(cmd)
}

// Exec is the function to exec any cmd on the session
func (p *PseudoTerminal) Exec(cmd string) ([]byte, error) {
	return p.Exec(cmd)
}

// Send is the function to send cmd to a pseudo terminal
func (p *PseudoTerminal) Send(cmd string) error {
	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.out == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.err == nil {
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

// UserCheck is the function to check current user is expected or not
func (p *PseudoTerminal) UserCheck(expected string) error {
	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.out == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.err == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	cmd := fmt.Sprintf("echo whoami; whoami\n")
	fmt.Fprintf(p.in, cmd)

	scanner := bufio.NewScanner(p.out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "whoami") {
			break
		}
	}

	for scanner.Scan() {
		user := scanner.Text()
		if user != expected {
			return fmt.Errorf("sudo su - %s was failed", expected)
		}
		return nil
	}
	return errors.New("Unexpected error happend at UserCheck func")
}

// SwitchUser is the function to switch user on pseudo terminal login sesion
func (p *PseudoTerminal) SwitchUser(user, passwd string) error {
	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.out == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.err == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		scanner := bufio.NewScanner(p.out)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "password") {
				wg.Done()
				return
			}
		}
	}()

	sudoCmd := fmt.Sprintf("sudo su - %s\n", user)
	fmt.Fprintf(p.in, sudoCmd)

	time.Sleep(1000 * time.Millisecond)
	inputPasswd := fmt.Sprintf("%s\n", passwd)
	fmt.Fprintf(p.in, inputPasswd)

	wg.Wait()
	return nil
}

// End is the function to end session as a pseudo terminal
func (p *PseudoTerminal) End() error {
	if p.in == nil {
		return errors.New("Psedudo stdin is not initialized")
	}

	if p.out == nil {
		return errors.New("Psedudo stdout is not initialized")
	}

	if p.err == nil {
		return errors.New("Psedudo stderr is not initialized")
	}

	fmt.Fprintf(p.in, "%s\n", "exit")
	return nil
}

// NewPseudoTerminal returns PseudoTerminal instance whose session is established by password.
func NewPseudoTerminal(w *Wrapper) (*PseudoTerminal, error) {
	inPipe, _ := w.session.StdinPipe()
	outPipe, _ := w.session.StdoutPipe()
	errPipe, _ := w.session.StderrPipe()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := w.session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	if err := w.session.Shell(); err != nil {
		return nil, fmt.Errorf("failed to start shell: %s", err)
	}

	return &PseudoTerminal{
		wrapper: w,
		in:      inPipe,
		out:     outPipe,
		err:     errPipe,
	}, nil
}
