package cmd

import (
	"fmt"

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
	execCmd.RunE = exec
	RootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, args []string) error {
	var dests []ssh.Dest
	var err error

	if destsYaml != "" {
		dests, err = ssh.NewDests(destsYaml, port)
		if err != nil {
			return err
		}
	}

	fmt.Println(dests)

	return nil
}
