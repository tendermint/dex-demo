package batch

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/types/errs"
)

const (
	QueryLatest = "latest"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryLatest:
			return queryLatest(path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown batch query endpoint")
		}
	}
}

func queryLatest(path []string, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, errs.ErrInvalidArgument("must specify a market ID")
	}

	marketID := sdk.NewUintFromString(path[0])
	res, sdkErr := keeper.LatestByMarket(marketID)
	if sdkErr != nil {
		if sdkErr.Code() == errs.CodeNotFound {
			return nil, nil
		}

		return nil, sdkErr
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, errs.ErrMarshalFailure("failed to marshal batch")
	}
	return b, nil
}
