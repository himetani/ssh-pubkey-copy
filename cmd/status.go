package cmd

import (
	"sync"

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

	client := ssh.PubKeyCopyClient{}
	rows := make([]ssh.Result, len(dests))

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, d := range dests {
		go func(i int, dest ssh.Dest) {
			defer wg.Done()

			session, err := ssh.NewPrivateKeySession(dest.Host, dest.Port, dest.User, privateKey)
			if err != nil {
				rows[i] = ssh.Result{Host: dest.Host, Port: dest.Port, User: dest.User, Err: err}
				return
			}
			defer session.Close()
			rows[i] = ssh.Result{Host: dest.Host, Port: dest.Port, User: dest.User, Err: client.Ping(session)}
			return
		}(i, d)
	}

	wg.Wait()
	table.Render(rows)

	return nil
}
