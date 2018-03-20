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

// NewDests return slice of Dest type
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
