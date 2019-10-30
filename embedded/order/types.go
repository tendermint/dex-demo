package order

import (
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Order struct {
	ID             store.EntityID     `json:"id"`
	Owner          sdk.AccAddress     `json:"owner"`
	MarketID       store.EntityID     `json:"market_id"`
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
	Start store.EntityID
	Owner sdk.AccAddress
}

type ListQueryResult struct {
	NextID store.EntityID `json:"next_id"`
	Orders []Order        `json:"orders"`
}
