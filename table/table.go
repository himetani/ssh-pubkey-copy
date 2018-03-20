package table

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Row is the struct that represents row in the table
type Row struct {
	User string
	Host string
	Port string
	Err  error
}

// Render is the function that render the table
func Render(rows []Row) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Destination", "Result"})
	table.SetRowLine(true)

	for _, r := range rows {
		if r.Err != nil {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[X] Not copied"})
		} else {
			table.Append([]string{fmt.Sprintf("%s@%s:%s", r.User, r.Host, r.Port), "[o] Already Copied"})
		}
	}
	table.Render()
}
