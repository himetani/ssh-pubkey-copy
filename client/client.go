package client

type Pinger interface {
	Ping()
}

type KeyClient struct {
}

func (k *KeyClient) Ping() []Result {
	return nil
}

type PasswordClient struct {
}

func (p *PasswordClient) Ping() []Result {
	return nil
}
