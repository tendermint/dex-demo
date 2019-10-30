package batch

import (
	"time"

	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Batch struct {
	BlockNumber   int64                     `json:"block_number"`
	BlockTime     time.Time                 `json:"block_time"`
	MarketID      store.EntityID            `json:"market_id"`
	ClearingPrice sdk.Uint                  `json:"clearing_price"`
	Bids          []matcheng.AggregatePrice `json:"bids"`
	Asks          []matcheng.AggregatePrice `json:"asks"`
}
