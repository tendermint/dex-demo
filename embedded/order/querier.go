package order

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryList = "list"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryList:
			return queryList(keeper, req.Data)
		default:
			return nil, sdk.ErrUnknownRequest("unknown embedded order request")
		}
	}
}

func queryList(keeper Keeper, reqB []byte) ([]byte, sdk.Error) {
	var req ListQueryRequest
	err := keeper.cdc.UnmarshalBinaryBare(reqB, &req)
	if err != nil {
		return nil, errs.ErrUnmarshalFailure("failed to unmarshal list query request")
	}

	orders := make([]Order, 0)
	var lastID store.EntityID
	iterCB := func(order Order) bool {
		orders = append(orders, order)
		lastID = order.ID
		return len(orders) < 50
	}

	if req.Owner.Empty() {
		if req.Start.IsDefined() {
			keeper.ReverseIteratorFrom(req.Start, iterCB)
		} else {
			keeper.ReverseIterator(iterCB)
		}
	} else {
		// TEMPORARY: can add support for richer querying with sqlite
		keeper.OrdersByOwner(req.Owner, iterCB)
	}

	if len(orders) < 50 {
		lastID = store.NewEntityID(0)
	}
	res := ListQueryResult{
		NextID: lastID.Dec(),
		Orders: orders,
	}
	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result")
	}
	return b, nil
}
