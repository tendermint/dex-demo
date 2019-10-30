package market

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/x/market/types"
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
			return nil, sdk.ErrUnknownRequest("unknown market query endpoint")
		}
	}
}

func queryList(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res := types.ListQueryResult{
		Markets: make([]types.NamedMarket, 0),
	}

	var retErr sdk.Error
	keeper.Iterator(ctx, func(mkt types.Market) bool {
		name, err := keeper.Pair(ctx, mkt.ID)
		if err != nil {
			retErr = err
			return false
		}

		res.Markets = append(res.Markets, types.NamedMarket{
			ID:           mkt.ID,
			BaseAssetID:  mkt.BaseAssetID,
			QuoteAssetID: mkt.QuoteAssetID,
			Name:         name,
		})
		return true
	})

	if retErr != nil {
		return nil, retErr
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		panic(err)
	}
	return b, nil
}
