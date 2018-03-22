package cmd

import (
	"sync"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/himetani/ssh-pubkey-copy/table"
	"github.com/spf13/cobra"
)

// bypassCmd represents the status command
var bypassCmd = &cobra.Command{
	Use:   "bypass",
	Short: "Execute copying the public key to target users of remote host bypassing the user",
	Long:  `Execute copying the public key to target users of remote host bypassing the user`,
}

func init() {
	bypassCmd.Flags().StringVar(&publicKeyPath, "key", "", "Public key to copy")

	bypassCmd.RunE = bypass
	RootCmd.AddCommand(bypassCmd)
}

func bypass(cmd *cobra.Command, args []string) error {
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
	var row table.Row

	var wg sync.WaitGroup
	wg.Add(1)

	dest := dests[0]
	go func() {
		defer wg.Done()
		terminal, err := ssh.NewPseudoTerminal(dest.Host, dest.Port, dest.User, password)
		if err != nil {
			row = table.Row{Host: dest.Host, Port: dest.Port, User: dest.User, Err: err}
			return
		}
		row = table.Row{Host: dest.Host, Port: dest.Port, User: dest.User, Err: client.BypassCopy(terminal, dest.User, content)}
		defer terminal.Close()
		return
	}()

	wg.Wait()

	return nil
}
