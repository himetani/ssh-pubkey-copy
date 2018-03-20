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

	client := ssh.PubKeyCopyClient{}
	rows := make([]table.Row, len(dests))

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, d := range dests {
		go func(i int, dest ssh.Dest) {
			defer wg.Done()

			session, err := ssh.NewPrivateKeySession(dest.Host, dest.Port, dest.User, privateKey)
			if err != nil {
				rows[i] = table.Row{Host: dest.Host, Port: dest.Port, User: dest.User, Err: err}
				return
			}
			defer session.Close()
			rows[i] = table.Row{Host: dest.Host, Port: dest.Port, User: dest.User, Err: client.Ping(session)}
			return
		}(i, d)
	}

	wg.Wait()
	table.Render(rows)

	return nil
}
