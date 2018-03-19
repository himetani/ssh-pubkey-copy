package ssh

import (
	"sync"
)

type Pinger interface {
	Ping()
}

type KeyClient struct {
	privateKey string
}

func (k *KeyClient) Ping(dests []Dest) []Result {
	results := []Result{}
	var wg sync.WaitGroup
	wg.Add(len(dests))

	for _, dest := range dests {
		dest := dest
		go func() {
			defer wg.Done()
			results = append(results, dest.ConnectWithPrivateKey(k.privateKey))
		}()
	}
	wg.Wait()

	return results
}

type PasswordClient struct {
	password string
}

func (p *PasswordClient) Ping() []Result {
	return nil
}

func NewKeyClient(privateKey string) *KeyClient {
	return &KeyClient{privateKey: privateKey}
}

func NewPasswordClient(privateKey string) *PasswordClient {
	return &PasswordClient{password: privateKey}
}
