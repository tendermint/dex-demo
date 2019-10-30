package order

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/x/order/types"
)

const (
	QueryList = "list"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryList:
			return queryList(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown order query endpoint")
		}
	}
}

func queryList(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res := types.ListQueryResult{
		Orders: make([]types.Order, 0),
	}

	keeper.Iterator(ctx, func(order types.Order) bool {
		res.Orders = append(res.Orders, order)
		return true
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		panic("could not marshal result")
	}
	return b, nil
}
