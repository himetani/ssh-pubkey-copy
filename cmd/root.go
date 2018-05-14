package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/olekukonko/tablewriter"
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
	RootCmd.PersistentFlags().StringVar(&privateKeyPath, "privateKey", "", "Selects a file from which the identity (private key) for public key authentication is read (default is $HOME/.ssh/id_rsa)")
	RootCmd.PersistentFlags().StringVarP(&destsYaml, "filename", "f", "", "input file")
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

// renderStatus is the function that render the status table
func renderStatus(rr []ssh.Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Status"})
	table.SetRowLine(true)

	for _, r := range rr {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[X] Not copied"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[o] Copied"})
		}
	}
	table.Render()
}

// renderStatus is the function that render the status table
func renderCopyResult(rr []ssh.Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Status", "Action", "Message"})
	table.SetRowLine(true)

	for _, r := range rr {
		if r.Err != nil {
			if _, ok := r.Err.(*ssh.SkipCopyError); ok {
				table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[o] Copied", "Skip", "Already copied"})
			} else {
				table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[X] Not copied", "Copy failed", r.Err.Error()})
			}
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[o] Copied", "Copy Success", ""})
		}
	}
	table.Render()
}
