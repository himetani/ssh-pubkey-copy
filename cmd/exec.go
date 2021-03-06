package cmd

import (
	"errors"

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
	execCmd.Flags().StringVar(&publicKeyPath, "key", "", "Selects a file from which the identity (public key) for public key authentication is read (default is $HOME/.ssh/id_rsa.pub")

	execCmd.RunE = exec
	RootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, args []string) error {
	if destsYaml == "" {
		return errors.New("Use -f, --filename option to specify input file")
	}

	var dests []ssh.Dest

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
		rr[i] = ssh.Copy(d.Host, d.Port, d.User, passwd, content, c)
	}

	for i, r := range rr {
		rows[i] = <-r
	}

	renderCopyResult(rows)

	return nil
}
