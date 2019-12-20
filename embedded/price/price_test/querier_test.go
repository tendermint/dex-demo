package price_test

import (
	"testing"
	"time"

	"github.com/tendermint/dex-demo/embedded/price"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/mockapp"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types"
)

func TestQuerier_Candles(t *testing.T) {
	testflags.UnitTest(t)
	app := mockapp.New(t)
	db := dbm.NewMemDB()
	keeper := price.NewKeeper(db, app.Cdc)
	mktID := sdk.OneUint()

	fills := []types.Fill{
		{
			sdk.OneUint(),
			mktID,
			testutil.RandAddr(),
			"DEX/ETH",
			matcheng.Bid,
			sdk.NewUint(100),
			sdk.ZeroUint(),
			1,
			100,
			sdk.NewUint(100),
		},
		{
			sdk.OneUint(),
			mktID,
			testutil.RandAddr(),
			"DEX/ETH",
			matcheng.Bid,
			sdk.NewUint(100),
			sdk.ZeroUint(),
			2,
			130,
			sdk.NewUint(90),
		},
		{
			sdk.OneUint(),
			mktID,
			testutil.RandAddr(),
			"DEX/ETH",
			matcheng.Bid,
			sdk.NewUint(100),
			sdk.ZeroUint(),
			3,
			160,
			sdk.NewUint(120),
		},
		{
			sdk.OneUint(),
			mktID,
			testutil.RandAddr(),
			"DEX/ETH",
			matcheng.Bid,
			sdk.NewUint(100),
			sdk.ZeroUint(),
			4,
			190,
			sdk.NewUint(140),
		},
	}

	for _, fill := range fills {
		keeper.OnFillEvent(fill)
	}
	querier := price.NewQuerier(keeper)

	t.Run("should support one minute candles", func(t *testing.T) {
		res := fetchResult(t, app.Ctx, querier, app.Cdc, 100, 190, price.CandleInterval1M)
		assert.Equal(t, 3, len(res.Candles))
		assertEqualCandleEntries(t, price.CandleEntry{
			Date:  time.Unix(60, 0),
			Open:  sdk.NewUint(100),
			Close: sdk.NewUint(100),
			High:  sdk.NewUint(100),
			Low:   sdk.NewUint(100),
		}, res.Candles[0])
		assertEqualCandleEntries(t, price.CandleEntry{
			Date:  time.Unix(120, 0),
			Open:  sdk.NewUint(90),
			Close: sdk.NewUint(120),
			High:  sdk.NewUint(120),
			Low:   sdk.NewUint(90),
		}, res.Candles[1])
	})
	t.Run("should support five minute candles", func(t *testing.T) {
		res := fetchResult(t, app.Ctx, querier, app.Cdc, 100, 190, price.CandleInterval5M)
		assert.Equal(t, 1, len(res.Candles))
		assertEqualCandleEntries(t, price.CandleEntry{
			Date:  time.Unix(0, 0),
			Open:  sdk.NewUint(100),
			Close: sdk.NewUint(140),
			High:  sdk.NewUint(140),
			Low:   sdk.NewUint(90),
		}, res.Candles[0])
	})
	t.Run("should support inexact start and end dates", func(t *testing.T) {
		res := fetchResult(t, app.Ctx, querier, app.Cdc, 101, 200, price.CandleInterval5M)
		assert.Equal(t, 1, len(res.Candles))
		assertEqualCandleEntries(t, price.CandleEntry{
			Date:  time.Unix(0, 0),
			Open:  sdk.NewUint(90),
			Close: sdk.NewUint(140),
			High:  sdk.NewUint(140),
			Low:   sdk.NewUint(90),
		}, res.Candles[0])
	})
}

func fetchResult(t *testing.T, ctx sdk.Context, querier sdk.Querier, cdc *amino.Codec, from int64, to int64, interval price.CandleInterval) price.CandleQueryResult {
	params := price.CandleQueryParams{
		From:     time.Unix(from, 0),
		To:       time.Unix(to, 0),
		Interval: interval,
	}
	paramsB := cdc.MustMarshalBinaryBare(params)
	req := abci.RequestQuery{
		Data: paramsB,
	}
	resJSON, err := querier(ctx, []string{"candles", "1"}, req)
	require.NoError(t, err)
	var res price.CandleQueryResult
	testutil.MustUnmarshalJSON(t, resJSON, &res)
	return res
}

func assertEqualCandleEntries(t *testing.T, expected price.CandleEntry, actual price.CandleEntry) {
	assert.Equal(t, expected.Date.Unix(), actual.Date.Unix())
	testutil.AssertEqualUints(t, expected.Open, actual.Open)
	testutil.AssertEqualUints(t, expected.Close, actual.Close)
	testutil.AssertEqualUints(t, expected.High, actual.High)
	testutil.AssertEqualUints(t, expected.Low, actual.Low)
}
