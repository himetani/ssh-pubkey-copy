package cmd

import (
	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/himetani/ssh-pubkey-copy/table"
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

	privateKey, err := ssh.NewPrivateKey(privateKeyPath)
	if err != nil {
		return err
	}

	if destsYaml != "" {
		dests, err = ssh.NewDests(destsYaml, port)
		if err != nil {
			return err
		}
	}

	rows := make([]ssh.Result, len(dests))

	rr := make([]<-chan ssh.Result, len(dests), len(dests))

	for i, d := range dests {
		rr[i] = ssh.IsCopy(d.Host, d.Port, d.User, privateKey)
	}

	for i, r := range rr {
		rows[i] = <-r
	}

	table.Render(rows)

	return nil
}
