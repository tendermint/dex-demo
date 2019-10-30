package types

import (
	"bytes"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type ListQueryResult struct {
	Orders []Order `json:"orders"`
}

func (l ListQueryResult) String() string {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	t.SetHeader([]string{
		"ID",
		"Owner",
		"MarketID",
		"Direction",
		"Price",
		"Quantity",
		"Time In Force",
		"Created Block",
	})

	for _, o := range l.Orders {
		t.Append([]string{
			o.ID.String(),
			o.Owner.String(),
			o.MarketID.String(),
			o.Direction.String(),
			o.Price.String(),
			o.Quantity.String(),
			strconv.FormatUint(uint64(o.TimeInForceBlocks), 10),
			strconv.Itoa(int(o.CreatedBlock)),
		})
	}
	t.Render()
	return string(buf.Bytes())
}
