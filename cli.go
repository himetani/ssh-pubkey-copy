package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
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
	flag.StringVar(&privateKey, "i", "", "Private key")
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
		results := targetsCmd(dests, privateKey)
		renderTable(results)
	}

	if subcmd == "copy" {
		fmt.Println("Not implemented yet")
	}

	return 0
}

func targetsCmd(dests []Dest, privateKey string) []Result {
	resultsChan := make(chan Result)
	results := []Result{}

	for _, dest := range dests {
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

func copyCmd() {

}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [copy|targets] [cfgFile]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func renderTable(results []Result) {
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
