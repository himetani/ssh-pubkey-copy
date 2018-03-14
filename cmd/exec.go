package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// execCmd represents the status command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute copying the public key to target users of remote host",
	Long:  `Execute copying the public key to target users of remote host`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exec cmd executed")
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
}
