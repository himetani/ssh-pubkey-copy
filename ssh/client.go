package ssh

import "fmt"

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
	if err := terminal.Start(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := terminal.SwitchUser(user, passwd); err != nil {
		fmt.Println(err)
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
