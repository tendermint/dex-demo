package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/matcheng"
)

type EventHandler interface {
	OnEvent(event interface{}) error
}

type Batch struct {
	BlockNumber   int64
	BlockTime     time.Time
	MarketID      sdk.Uint
	ClearingPrice sdk.Uint
	Bids          []matcheng.AggregatePrice
	Asks          []matcheng.AggregatePrice
}

type Fill struct {
	OrderID     sdk.Uint
	MarketID    sdk.Uint
	Owner       sdk.AccAddress
	Pair        string
	Direction   matcheng.Direction
	QtyFilled   sdk.Uint
	QtyUnfilled sdk.Uint
	BlockNumber int64
	BlockTime   int64
	Price       sdk.Uint
}

type OrderCreated struct {
	ID                sdk.Uint
	Owner             sdk.AccAddress
	MarketID          sdk.Uint
	Direction         matcheng.Direction
	Price             sdk.Uint
	Quantity          sdk.Uint
	TimeInForceBlocks uint16
	CreatedBlock      int64
}

type OrderCancelled struct {
	OrderID sdk.Uint
}

type BurnCreated struct {
	ID          sdk.Uint
	AssetID     sdk.Uint
	BlockNumber int64
	Burner      sdk.AccAddress
	Beneficiary []byte
	Quantity    sdk.Uint
}
