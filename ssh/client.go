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
	resultsChan := make(chan Result)
	results := []Result{}

	for _, dest := range dests {
		dest := dest
		go dest.ConnectWithPrivateKey(resultsChan, k.privateKey)
	}

	var wg sync.WaitGroup
	wg.Add(len(dests))

	go func() {
		for r := range resultsChan {
			results = append(results, r)
			wg.Done()
		}
	}()

	wg.Wait()
	close(resultsChan)

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
