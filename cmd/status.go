package cmd

import (
	"path/filepath"
	"sync"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/himetani/ssh-pubkey-copy/table"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the connectivity status of target users of remote host",
	Long:  `Show the connectivity status of target users of remote host`,
}

func init() {
	statusCmd.RunE = status
	RootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) error {
	var dests []ssh.Dest
	var err error

	if privateKey == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		privateKey = filepath.Join(home, ".ssh", "id_rsa")
	}

	if destsYaml != "" {
		dests, err = ssh.NewDests(destsYaml, port)
		if err != nil {
			return err
		}
	}

	client := ssh.NewKeyClient(privateKey)
	results := []ssh.Result{}
	var wg sync.WaitGroup
	wg.Add(len(dests))

	for _, dest := range dests {
		dest := dest
		go func() {
			defer wg.Done()
			results = append(results, client.Ping(dest))
			return
		}()
	}
	wg.Wait()
	table.Render(results)

	return nil
}
