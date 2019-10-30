package market_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/mockapp"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper(t *testing.T) {
	testflags.UnitTest(t)
	app := mockapp.New(t)
	asset1, err := app.AssetKeeper.Create(app.Ctx, "test asset", "TST1", testutil.RandAddr(), sdk.NewUint(1000000))
	require.NoError(t, err)
	asset2, err := app.AssetKeeper.Create(app.Ctx, "test asset", "TST2", testutil.RandAddr(), sdk.NewUint(1000000))
	require.NoError(t, err)
	mkt := app.MarketKeeper.Create(app.Ctx, asset1.ID, asset2.ID)
	expMkt := types.Market{
		ID:           store.NewEntityID(1),
		BaseAssetID:  asset1.ID,
		QuoteAssetID: asset2.ID,
	}

	assert.EqualValues(t, expMkt, mkt)

	retMkt, err := app.MarketKeeper.Get(app.Ctx, mkt.ID)
	require.NoError(t, err)
	assert.EqualValues(t, expMkt, retMkt)

	assert.True(t, app.MarketKeeper.Has(app.Ctx, mkt.ID))

	pair, err := app.MarketKeeper.Pair(app.Ctx, mkt.ID)
	require.NoError(t, err)
	assert.Equal(t, "TST1/TST2", pair)
}
