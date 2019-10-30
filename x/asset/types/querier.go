package types

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type ListQueryResult struct {
	Assets []Asset `json:"assets"`
}

func (l ListQueryResult) String() string {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	t.SetHeader([]string{
		"ID",
		"Name",
		"Symbol",
		"Owner",
		"Circulating Supply",
		"Total Supply",
	})

	for _, a := range l.Assets {
		t.Append([]string{
			a.ID.String(),
			a.Name,
			a.Symbol,
			a.Owner.String(),
			a.CirculatingSupply.String(),
			a.TotalSupply.String(),
		})
	}

	t.Render()
	return string(buf.Bytes())
}
