package main

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
	return dests, nil
}

func (d *Dest) ConnectWithPrivateKey(resultChan chan<- Result, privateKey string) {
	session, err := NewSession(d.Host, d.Port, d.User, privateKey)
	if err != nil {
		resultChan <- Result{Dest: d, Err: err}
		return
	}
	defer session.Close()

	bytes, err := session.Connect()
	fmt.Println(string(bytes))
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

type Result struct {
	*Dest
	Err error
}
