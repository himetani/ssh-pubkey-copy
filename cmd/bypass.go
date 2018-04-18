package cmd

import (
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

	bypassUser := args[0]

	content, err := ssh.NewPublicKeyContent(publicKeyPath)
	if err != nil {
		return err
	}

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

	passwd, err := getPassword()
	if err != nil {
		return err
	}

	rows := make([]ssh.Result, len(dests))

	rr := make([]<-chan ssh.Result, len(dests), len(dests))

	for i, d := range dests {
		c := ssh.IsCopy(d.Host, d.Port, d.User, privateKey)
		rr[i] = ssh.BypassCopy(d.Host, d.Port, d.User, passwd, bypassUser, content, c)
	}

	for i, r := range rr {
		rows[i] = <-r
	}

	renderCopyResult(rows)

	return nil
}
