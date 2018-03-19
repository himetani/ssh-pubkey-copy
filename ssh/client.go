package ssh

import (
	"sync"
)

type Pinger interface {
	Ping()
}

type Client interface {
	Pinger
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

			session, err := NewPrivateKeySession(dest.Host, dest.Port, dest.User, k.privateKey)
			if err != nil {
				results = append(results, Result{Dest: &dest, Err: err})
				return
			}
			defer session.Close()

			_, err = session.Connect()
			if err != nil {
				results = append(results, Result{Dest: &dest, Err: err})
				return
			}

			results = append(results, Result{Dest: &dest, Err: nil})
		}()
	}

	wg.Wait()

	return results
}

type PasswordClient struct {
	password string
}

func (p *PasswordClient) Ping(dests []Dest) []Result {
	return nil
}

func NewKeyClient(privateKey string) *KeyClient {
	return &KeyClient{privateKey: privateKey}
}

func NewPasswordClient(password string) *PasswordClient {
	return &PasswordClient{password: password}
}
