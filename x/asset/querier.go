package asset

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/x/asset/types"
)

const (
	QueryList = "list"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryList:
			return queryList(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryList(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	symbol := strings.ToUpper(path[0])

	res := types.ListQueryResult{
		Assets: make([]types.Asset, 0),
	}

	keeper.Iterator(ctx, func(asset types.Asset) bool {
		if symbol == "" || strings.Contains(asset.Symbol, symbol) {
			res.Assets = append(res.Assets, asset)
		}
		return true
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		panic("could not marshal result")
	}

	return b, nil
}
