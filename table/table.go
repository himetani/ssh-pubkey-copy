package table

import (
	"fmt"
	"os"

	"github.com/himetani/ssh-pubkey-copy/client"
	"github.com/olekukonko/tablewriter"
)

func Render(results []client.Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetRowLine(true)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.FgHiRedColor})

	for _, r := range results {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.Dest.User, r.Dest.Host, r.Dest.Port), "🙅♀️Invalid Host"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.Dest.User, r.Dest.Host, r.Dest.Port), "🙆PubKey is Copied"})
		}
	}
	table.Render()
}
