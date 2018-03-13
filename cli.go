package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

type CLI struct {
}

func (c *CLI) Run() int {
	var (
		help       bool
		quiet      bool
		privateKey string
		port       string
	)

	const (
		defaultPort = "22"
	)

	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&quiet, "q", false, "Don't show the INFO log")
	flag.StringVar(&privateKey, "i", "", "Private key (Default: $HOME/.ssh/id_rsa.pub")
	flag.StringVar(&port, "p", "", "Port")
	flag.Parse()

	if help {
		showUsage()
		return 0
	}

	args := flag.Args()
	if len(args) != 2 {
		showUsage()
		return 1
	}

	if args[0] != "copy" && args[0] != "targets" {
		showUsage()
		return 1
	}

	subcmd := args[0]

	if port == "" {
		port = defaultPort
	}

	dests, err := NewDests(args[1], port)
	if err != nil {
		fmt.Printf(err.Error())
		return 1
	}

	if privateKey == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return 1
		}
		privateKey = filepath.Join(home, ".ssh", "id_rsa")
	}

	if subcmd == "targets" {
		targetsCmd(dests, privateKey)
	}

	if subcmd == "copy" {
		pubkey := privateKey + ".pub"
		if err := copyCmd(dests, pubkey); err != nil {
			fmt.Println(err.Error())
			return 1
		}
	}

	return 0
}

func targetsCmd(dests []Dest, privateKey string) {
	results := connectWithKey(dests, privateKey)
	renderTable(results)
}

func copyCmd(dests []Dest, pubkey string) error {
	password, err := getPassword()
	if err != nil {
		return err
	}

	pubKeyContent, err := getPubKeyContent(pubkey)
	if err != nil {
		return err
	}

	results := connectWithPassword(dests, password, pubKeyContent)
	renderTable(results)

	return nil
}

func getPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("")

	return string(password), nil
}

func getPubKeyContent(pubKeyPath string) (string, error) {
	content, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func connectWithKey(dests []Dest, privateKey string) []Result {
	resultsChan := make(chan Result)
	results := []Result{}

	for _, dest := range dests {
		dest := dest
		go dest.ConnectWithPrivateKey(resultsChan, privateKey)
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

func connectWithPassword(dests []Dest, password, pubKeyContent string) []Result {
	resultsChan := make(chan Result)
	results := []Result{}

	for _, dest := range dests {
		dest := dest
		go dest.ExecWithPassword(resultsChan, password, pubKeyContent)
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

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [copy|targets] [cfgFile]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func renderTable(results []Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetRowLine(true)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})

	for _, r := range results {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.Dest.User, r.Dest.Host, r.Dest.Port), "🙅♀️Invalid Host"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s", r.Dest.User, r.Dest.Host), "🙆PubKey is Copied"})
		}
	}
	table.Render()
}
