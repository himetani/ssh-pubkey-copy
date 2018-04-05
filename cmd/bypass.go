package cmd

import (
	"fmt"
	"sync"

	"github.com/himetani/ssh-pubkey-copy/ssh"
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

	bypass := args[0]

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

	passwd, err := getPassword()
	if err != nil {
		return err
	}

	client := ssh.PubKeyCopyClient{}
	var row ssh.Result

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for _, dest := range dests {
		dest := dest
		fmt.Println(dest.Host)
		go func() {
			defer wg.Done()
			wrapper, err := ssh.NewPasswordSession(dest.Host, dest.Port, bypass, passwd)
			terminal, err := ssh.NewPseudoTerminal(wrapper)
			if err != nil {
				row = ssh.Result{Host: dest.Host, Port: dest.Port, User: dest.User, Err: err}
				return
			}
			row = ssh.Result{Host: dest.Host, Port: dest.Port, User: dest.User, Err: client.BypassCopy(terminal, dest.User, passwd, content)}
			defer terminal.Close()
			return
		}()
	}

	wg.Wait()

	return nil
}
