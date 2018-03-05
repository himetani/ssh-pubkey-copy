package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh"
)

var help bool
var quiet bool
var privateKey string
var port string

const (
	defaultPort = "22"
)

type Dest struct {
	Host string `json:"host"`
	User string `json:"user"`
}

type Result struct {
	Dest
	Err error
}

func init() {
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&quiet, "q", false, "Don't show the INFO log")
	flag.StringVar(&privateKey, "i", "", "Private key")
	flag.StringVar(&port, "p", "", "Port")
	flag.Parse()
}

func main() {
	if help {
		showUsage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		showUsage()
		os.Exit(1)
	}

	if args[0] != "copy" && args[0] != "targets" {
		showUsage()
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	var dests []Dest
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dests); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if privateKey == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		privateKey = filepath.Join(home, ".ssh", "id_rsa")
	}

	if port == "" {
		port = defaultPort
	}

	resultsChan := make(chan Result)
	results := []Result{}

	for _, dest := range dests {
		go call(resultsChan, dest, port)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})

	for _, r := range results {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s", r.Dest.User, r.Dest.Host), "🙅♀️Invalid Host"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s", r.Dest.User, r.Dest.Host), "🙆PubKey is Copied"})
		}
	}
	table.Render()
}

func call(results chan<- Result, dest Dest, port string) {
	session, err := NewSession(dest.Host, port, dest.User, privateKey)
	if err != nil {
		results <- Result{Dest: dest, Err: err}
		return
	}
	defer session.Close()

	bytes, err := session.Connect()
	fmt.Println(string(bytes))
	if err != nil {
		results <- Result{Dest: dest, Err: err}
		return
	}

	if err != nil {
		results <- Result{Dest: dest, Err: err}
		return
	}

	results <- Result{Dest: dest, Err: nil}
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [copy|targets] [cfgFile]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

// Session is struct representing ssh Session
type Session struct {
	config    *ssh.ClientConfig
	conn      *ssh.Client
	session   *ssh.Session
	StdinPipe io.WriteCloser
}

// NewSession returns new Session instance
func NewSession(ip, port, user, privateKey string) (*Session, error) {
	buf, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: 1 * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return &Session{
		config:  config,
		conn:    conn,
		session: session,
	}, nil
}

// Close close the session & connection
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}

// Get is func that get file contents
func (s *Session) Connect() ([]byte, error) {
	cmd := fmt.Sprintf("echo 'connect'\n")
	return s.session.Output(cmd)
}
