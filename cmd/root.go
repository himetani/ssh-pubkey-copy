package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var port string
var privateKeyPath string
var publicKeyPath string
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
	RootCmd.PersistentFlags().StringVarP(&privateKeyPath, "identity-file", "i", "", "Selects a file from which the identity (private key) for public key authentication is read (default is $HOME/.ssh/id_rsa.pub)")
	RootCmd.PersistentFlags().StringVarP(&destsYaml, "dests-file", "d", "", "dests.yml")
}

func getPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("")

	return string(password), nil
}
