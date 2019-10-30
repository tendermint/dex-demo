package types

import (
	"bytes"

	"github.com/olekukonko/tablewriter"

	"github.com/tendermint/dex-demo/types/store"
)

type NamedMarket struct {
	ID           store.EntityID
	BaseAssetID  store.EntityID
	QuoteAssetID store.EntityID
	Name         string
}

type ListQueryResult struct {
	Markets []NamedMarket `json:"markets"`
}

func (l ListQueryResult) String() string {
	var buf bytes.Buffer
	t := tablewriter.NewWriter(&buf)
	t.SetHeader([]string{
		"ID",
		"Name",
		"Base Asset ID",
		"Quote Asset ID",
	})

	for _, m := range l.Markets {
		t.Append([]string{
			m.ID.String(),
			m.Name,
			m.BaseAssetID.String(),
			m.QuoteAssetID.String(),
		})
	}

	t.Render()
	return string(buf.Bytes())
}
