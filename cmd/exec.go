package cmd

import (
	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/himetani/ssh-pubkey-copy/table"
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

	passwd, err := getPassword()
	if err != nil {
		return err
	}

	rows := make([]ssh.Result, len(dests))

	rr := make([]<-chan ssh.Result, len(dests), len(dests))

	for i, d := range dests {
		rr[i] = ssh.Copy(d.Host, d.Port, d.User, passwd, content)
	}

	for i, r := range rr {
		rows[i] = <-r
	}

	table.Render(rows)

	return nil
}
