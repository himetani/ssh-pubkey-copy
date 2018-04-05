package table

import (
	"fmt"
	"os"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/olekukonko/tablewriter"
)

// Render is the function that render the table
func Render(rows []ssh.Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetRowLine(true)

	for _, r := range rows {
		if r.Err != nil {
			if serr, ok := r.Err.(*ssh.ConnectionError); ok {
				table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), serr.Error()})
			} else {
				table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[X] Not copied"})
			}
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[o] Already Copied"})
		}
	}
	table.Render()
}
