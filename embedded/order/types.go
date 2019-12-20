package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/matcheng"
)

type Order struct {
	ID             sdk.Uint           `json:"id"`
	Owner          sdk.AccAddress     `json:"owner"`
	MarketID       sdk.Uint           `json:"market_id"`
	Direction      matcheng.Direction `json:"direction"`
	Price          sdk.Uint           `json:"price"`
	Quantity       sdk.Uint           `json:"quantity"`
	Status         string             `json:"status"`
	Type           string             `json:"type"`
	TimeInForce    uint16             `json:"time_in_force"`
	QuantityFilled sdk.Uint           `json:"quantity_filled"`
	CreatedBlock   int64              `json:"created_block"`
}

type ListQueryRequest struct {
	Start sdk.Uint
	Owner sdk.AccAddress
}

type ListQueryResult struct {
	NextID sdk.Uint `json:"next_id"`
	Orders []Order  `json:"orders"`
}
