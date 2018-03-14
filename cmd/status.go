package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/himetani/ssh-pubkey-copy/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the connectivity status of target users of remote host",
	Long:  `Show the connectivity status of target users of remote host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status cmd called")
	},
}

func init() {
	statusCmd.RunE = status
	RootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) error {
	var dests []client.Dest
	var err error

	if privateKey == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		privateKey = filepath.Join(home, ".ssh", "id_rsa")
	}

	if destsYaml != "" {
		dests, err = client.NewDests(destsYaml, port)
		if err != nil {
			return err
		}
	}

	results := connectWithKey(dests, privateKey)
	renderTable(results)

	return nil
}

func connectWithKey(dests []client.Dest, privateKey string) []client.Result {
	resultsChan := make(chan client.Result)
	results := []client.Result{}

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

func renderTable(results []client.Result) {
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
