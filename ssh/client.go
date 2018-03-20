package ssh

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
