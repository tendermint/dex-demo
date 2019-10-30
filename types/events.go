package types

import (
	"time"

	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EventHandler interface {
	OnEvent(event interface{}) error
}

type Batch struct {
	BlockNumber   int64
	BlockTime     time.Time
	MarketID      store.EntityID
	ClearingPrice sdk.Uint
	Bids          []matcheng.AggregatePrice
	Asks          []matcheng.AggregatePrice
}

type Fill struct {
	OrderID     store.EntityID
	MarketID    store.EntityID
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
	ID                store.EntityID
	Owner             sdk.AccAddress
	MarketID          store.EntityID
	Direction         matcheng.Direction
	Price             sdk.Uint
	Quantity          sdk.Uint
	TimeInForceBlocks uint16
	CreatedBlock      int64
}

type OrderCancelled struct {
	OrderID store.EntityID
}

type BurnCreated struct {
	ID          store.EntityID
	AssetID     store.EntityID
	BlockNumber int64
	Burner      sdk.AccAddress
	Beneficiary []byte
	Quantity    sdk.Uint
}
