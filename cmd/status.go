package cmd

import (
	"path/filepath"

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
	results := client.Ping(dests)
	table.Render(results)

	return nil
}
