package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var port string
var privateKey string
var destsYaml string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ssh-pubkey-copy",
	Short: "Copy ssh public key to target user of remote host",
	Long:  `Copy ssh public key to target user of remote host`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&port, "port", "22", "Port to connect to on the remote host")
	RootCmd.PersistentFlags().StringVar(&privateKey, "privateKey", "", "Selects a file from which the identity (private key) for public key authentication is read (default is $HOME/.ssh/id_rsa.pub)")
	RootCmd.PersistentFlags().StringVar(&destsYaml, "dests", "", "dests.yml")
}
