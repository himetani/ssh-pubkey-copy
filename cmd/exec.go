package cmd

import (
	"sync"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/spf13/cobra"
)

// execCmd represents the status command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute copying the public key to target users of remote host",
	Long:  `Execute copying the public key to target users of remote host`,
}

func init() {
	execCmd.Flags().StringVar(&publicKeyPath, "key", "", "Public key to copy")

	execCmd.RunE = exec
	RootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, args []string) error {
	var dests []ssh.Dest

	content, err := ssh.NewPublicKeyContent(publicKeyPath)
	if err != nil {
		return err
	}

	if destsYaml != "" {
		dests, err = ssh.NewDests(destsYaml, port)
		if err != nil {
			return err
		}
	}

	password, err := getPassword()
	if err != nil {
		return err
	}

	client := ssh.PubKeyCopyClient{}
	rows := make([]ssh.Result, len(dests))

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, d := range dests {
		go func(i int, dest ssh.Dest) {
			defer wg.Done()

			session, err := ssh.NewPasswordSession(dest.Host, dest.Port, dest.User, password)
			if err != nil {
				rows[i] = ssh.Result{Host: dest.Host, Port: dest.Port, User: dest.User, Err: err}
				return
			}
			defer session.Close()

			rows[i] = ssh.Result{
				Host: dest.Host,
				Port: dest.Port,
				User: dest.User,
				Err:  client.Copy(session, content),
			}
			return
		}(i, d)
	}

	wg.Wait()

	return nil
}
