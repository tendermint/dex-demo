package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/matcheng"
)

const (
	ModuleName = "order"
	RouterKey  = ModuleName
	StoreKey   = ModuleName
)

const MaxTimeInForce = 600

type Order struct {
	ID                sdk.Uint           `json:"id"`
	Owner             sdk.AccAddress     `json:"owner"`
	MarketID          sdk.Uint           `json:"market"`
	Direction         matcheng.Direction `json:"direction"`
	Price             sdk.Uint           `json:"price"`
	Quantity          sdk.Uint           `json:"quantity"`
	TimeInForceBlocks uint16             `json:"time_in_force_blocks"`
	CreatedBlock      int64              `json:"created_block"`
}

func New(owner sdk.AccAddress, marketID sdk.Uint, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16, created int64) Order {
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
