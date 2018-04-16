package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// IsCopy is func to check that ssh's public key is already copied or not
func IsCopy(ip, port, user string, privateKey ssh.Signer) chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		session, err := NewPrivateKeySession(ip, port, user, privateKey)
		if err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}
		defer session.Close()
		out <- Result{Host: ip, Port: port, User: user, Err: nil}
		return
	}()
	return out
}

// Copy is func to copy ssh's public key to target user
func Copy(ip, port, user, passwd, content string, in <-chan Result) chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)

		r := <-in
		if r.Err == nil {
			out <- Result{Host: ip, Port: port, User: user, Err: nil}
			return
		}

		session, err := NewPasswordSession(ip, port, user, passwd)
		if err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}
		defer session.Close()

		cmd := fmt.Sprintf("mkdir -p $HOME/.ssh; chmod 755 $HOME/.ssh; touch $HOME/.ssh/authorized_keys; chmod 600 $HOME/.ssh/authorized_keys; echo '%s'>> $HOME/.ssh/authorized_keys", content)

		_, err = session.Exec(cmd)

		out <- Result{
			Host: ip,
			Port: port,
			User: user,
			Err:  err,
		}
		return
	}()
	return out
}

// BypassCopy is func that copy ssh's public key to target user via bypass user ssh session
func BypassCopy(ip, port, user, passwd, bypassUser, content string, in <-chan Result) chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)

		r := <-in
		if r.Err == nil {
			out <- Result{Host: ip, Port: port, User: user, Err: nil}
			return
		}

		wrapper, err := NewPasswordSession(ip, port, bypassUser, passwd)
		if err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}

		terminal, err := NewPseudoTerminal(wrapper)
		if err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}

		if err := terminal.SwitchUser(user, passwd); err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}

		cmd := fmt.Sprintf("mkdir -p $HOME/.ssh; chmod 755 $HOME/.ssh;touch $HOME/.ssh/authorized_keys;chmod 600 $HOME/.ssh/authorized_keys;echo '%s'>>$HOME/.ssh/authorized_keys", content)
		if err := terminal.Send(cmd); err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}

		if err := terminal.End(); err != nil {
			out <- Result{Host: ip, Port: port, User: user, Err: err}
			return
		}

		out <- Result{
			Host: ip,
			Port: port,
			User: user,
			Err:  nil,
		}
		return
	}()
	return out
}
