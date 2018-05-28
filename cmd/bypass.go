package cmd

import (
	"errors"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/spf13/cobra"
)

// bypassCmd represents the status command
var bypassCmd = &cobra.Command{
	Use:   "bypass [bypassUserName]",
	Short: "Execute copying the public key to target users of remote host bypassing the user",
	Long:  `Execute copying the public key to target users of remote host bypassing the user`,
}

func init() {
	bypassCmd.Flags().StringVar(&publicKeyPath, "publicKey", "", "Selects a file from which the identity (public key) for public key authentication is read (default is $HOME/.ssh/id_rsa.pub")

	bypassCmd.RunE = bypass
	RootCmd.AddCommand(bypassCmd)
}

func bypass(cmd *cobra.Command, args []string) error {
	if destsYaml == "" {
		return errors.New("Use -f, --filename option to specify input file")
	}

	if len(args) != 1 {
		return errors.New("Invalid arguments")
	}

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
		tmpCh := ssh.BypassCopy(d.Host, d.Port, d.User, passwd, bypassUser, content, c)
		rr[i] = ssh.IsCopyWithChan(tmpCh, privateKey)
	}

	for i, r := range rr {
		rows[i] = <-r
	}

	renderCopyResult(rows)

	return nil
}
