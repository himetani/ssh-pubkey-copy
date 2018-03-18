package ssh

import (
	"encoding/json"
	"fmt"
	"os"
)

type Dest struct {
	Host string `json:"host"`
	Port string
	User string `json:"user"`
}

func NewDests(cfgPath, port string) ([]Dest, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil, err
	}

	var dests []Dest
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil, err
	}

	for i, _ := range dests {
		dests[i].Port = port
	}
	return dests, nil
}

func (d *Dest) ConnectWithPrivateKey(resultChan chan<- Result, privateKey string) {
	session, err := NewPrivateKeySession(d.Host, d.Port, d.User, privateKey)
	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}

	defer session.Close()

	_, err = session.Connect()
	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}

	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}

	resultChan <- Result{Dest: d, Err: nil}
	return
}

func (d *Dest) ExecWithPassword(resultChan chan<- Result, password, pubkey string) {
	session, err := NewPasswordSession(d.Host, d.Port, "vagrant", password)
	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}
	defer session.Close()

	cmd := fmt.Sprintf("echo %s | sudo -S su %s; mkdir -p .ssh/; echo \"%s\" >> .ssh/authorized_keys", password, d.User, pubkey)

	_, err = session.Exec(cmd)
	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}

	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}

	resultChan <- Result{Dest: d, Err: nil}
}
