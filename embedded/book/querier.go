package book

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/dex-demo/embedded/order"
	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGet = "get"
)

func NewQuerier(keeper order.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown spread query endpoint")
		}
	}
}

func queryGet(path []string, keeper order.Keeper) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, errs.ErrInvalidArgument("must specify a market ID")
	}

	mktId := store.NewEntityIDFromString(path[0])
	res := keeper.OpenOrdersByMarket(mktId)
	b, err := codec.MarshalJSONIndent(codec.New(), res)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result")
	}
	return b, nil
}
