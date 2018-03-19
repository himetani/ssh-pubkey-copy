package table

import (
	"fmt"
	"os"

	"github.com/himetani/ssh-pubkey-copy/ssh"
	"github.com/olekukonko/tablewriter"
)

func Render(results []ssh.Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetRowLine(true)
	//table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})

	for _, r := range results {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.Dest.User, r.Dest.Host, r.Dest.Port), "[X] Not copied"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.Dest.User, r.Dest.Host, r.Dest.Port), "[o] Already Copied"})
		}
	}
	table.Render()
}
