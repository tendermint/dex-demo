package types

import (
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "order"
	RouterKey  = ModuleName
	StoreKey   = ModuleName
)

const MaxTimeInForce = 600

type Order struct {
	ID                store.EntityID     `json:"id"`
	Owner             sdk.AccAddress     `json:"owner"`
	MarketID          store.EntityID     `json:"market"`
	Direction         matcheng.Direction `json:"direction"`
	Price             sdk.Uint           `json:"price"`
	Quantity          sdk.Uint           `json:"quantity"`
	TimeInForceBlocks uint16             `json:"time_in_force_blocks"`
	CreatedBlock      int64              `json:"created_block"`
}

func New(owner sdk.AccAddress, marketID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16, created int64) Order {
	return Order{
		Owner:             owner,
		MarketID:          marketID,
		Direction:         direction,
		Price:             price,
		Quantity:          quantity,
		TimeInForceBlocks: tif,
		CreatedBlock:      created,
	}
}
