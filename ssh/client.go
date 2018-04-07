package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Pinger interface {
	Ping()
}

type Client interface {
	Pinger
}

type PubKeyCopyClient struct{}

// Ping is func to send ping using session
func (p *PubKeyCopyClient) Ping(session Session) error {
	_, err := session.Connect()
	return err
}

// Copy is func to copy the private key using session
func (p *PubKeyCopyClient) Copy(session Session, publicKey string) error {
	cmd := fmt.Sprintf("mkdir -p $HOME/.ssh; chmod 755 $HOME/.ssh; touch $HOME/.ssh/authorized_keys; chmod 600 $HOME/.ssh/authorized_keys; echo '%s'>> $HOME/.ssh/authorized_keys", publicKey)
	_, err := session.Exec(cmd)
	return err
}

// Copy is func to copy the private key using session
func (p *PubKeyCopyClient) BypassCopy(terminal Terminal, user, passwd, publicKey string) error {
	if err := terminal.SwitchUser(user, passwd); err != nil {
		return err
	}

	cmd := fmt.Sprintf("mkdir -p $HOME/.ssh; chmod 755 $HOME/.ssh;touch $HOME/.ssh/authorized_keys;chmod 600 $HOME/.ssh/authorized_keys;echo '%s'>>$HOME/.ssh/authorized_keys", publicKey)
	if err := terminal.Send(cmd); err != nil {
		fmt.Println(err)
		return err
	}

	if err := terminal.End(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

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

func Copy(ip, port, user string, privateKey ssh.Signer, resultChan chan<- Result) {
	go func() {
		/*
			defer close(resultChan)

			session, err := NewPrivateKeySession(ip, port, user, privateKey)
			if err != nil {
				resultChan <- Result{Host: ip, Port: port, User: user, Err: err}
				return
			}
			defer session.Close()
			resultChan <- Result{Host: ip, Port: port, User: user, Err: nil}
			return
		*/
	}()
}
