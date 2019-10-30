package execution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/mockapp"
	"github.com/tendermint/dex-demo/testutil/testflags"
	uexstore "github.com/tendermint/dex-demo/types/store"
)

func TestKeeper_ExecuteAndCancelExpired(t *testing.T) {
	testflags.UnitTest(t)
	app := mockapp.New(t)
	owner := testutil.RandAddr()
	buyer := testutil.RandAddr()
	seller := testutil.RandAddr()

	asset1, err := app.AssetKeeper.Create(app.Ctx, "test asset", "TST1", owner, testutil.ToBaseUnits(1000000))
	require.NoError(t, err)
	asset2, err := app.AssetKeeper.Create(app.Ctx, "test asset", "TST2", owner, testutil.ToBaseUnits(1000000))
	require.NoError(t, err)
	require.NoError(t, app.AssetKeeper.Mint(app.Ctx, asset1.ID, testutil.ToBaseUnits(1000000)))
	require.NoError(t, app.AssetKeeper.Mint(app.Ctx, asset2.ID, testutil.ToBaseUnits(1000000)))
	require.NoError(t, app.AssetKeeper.Transfer(app.Ctx, asset1.ID, owner, buyer, testutil.ToBaseUnits(10000)))
	require.NoError(t, app.AssetKeeper.Transfer(app.Ctx, asset2.ID, owner, buyer, testutil.ToBaseUnits(10000)))
	require.NoError(t, app.AssetKeeper.Transfer(app.Ctx, asset1.ID, owner, seller, testutil.ToBaseUnits(10000)))
	require.NoError(t, app.AssetKeeper.Transfer(app.Ctx, asset2.ID, owner, seller, testutil.ToBaseUnits(10000)))
	mkt := app.MarketKeeper.Create(app.Ctx, asset1.ID, asset2.ID)

	_, err = app.OrderKeeper.Post(app.Ctx, buyer, mkt.ID, matcheng.Bid, testutil.ToBaseUnits(1), testutil.ToBaseUnits(1), 100)
	require.NoError(t, err)

	ctx := app.Ctx.WithBlockHeight(602)
	bids := [][2]uint64{
		{1, 10},
		{2, 10},
		{3, 10},
	}
	asks := [][2]uint64{
		{2, 10},
		{3, 10},
		{4, 10},
	}
	for _, bid := range bids {
		_, err = app.OrderKeeper.Post(ctx, buyer, mkt.ID, matcheng.Bid, testutil.ToBaseUnits(bid[0]), testutil.ToBaseUnits(bid[1]), 100)
		require.NoError(t, err)
	}
	for _, ask := range asks {
		_, err = app.OrderKeeper.Post(ctx, seller, mkt.ID, matcheng.Ask, testutil.ToBaseUnits(ask[0]), testutil.ToBaseUnits(ask[1]), 100)
		require.NoError(t, err)
	}
	require.NoError(t, app.ExecutionKeeper.ExecuteAndCancelExpired(ctx))
	t.Run("should expire orders out of TIF", func(t *testing.T) {
		assert.False(t, app.OrderKeeper.Has(ctx, uexstore.NewEntityID(1)))
	})
	t.Run("should update quantities of partially filled orders", func(t *testing.T) {
		ord3, err := app.OrderKeeper.Get(ctx, uexstore.NewEntityID(3))
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(5), ord3.Quantity)
		ord4, err := app.OrderKeeper.Get(ctx, uexstore.NewEntityID(4))
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(5), ord4.Quantity)
	})

	// perform next round of cancellation after since orders are
	// deleted on cancellation
	ctx = app.Ctx.WithBlockHeight(703)
	require.NoError(t, app.ExecutionKeeper.ExecuteAndCancelExpired(ctx))

	t.Run("should delete completely filled orders", func(t *testing.T) {
		assert.False(t, app.OrderKeeper.Has(ctx, uexstore.NewEntityID(5)))
	})
	t.Run("all executed orders should exchange coins", func(t *testing.T) {
		// seller should have 9990 asset 1, because two orders were
		// partially executed (for 5 each), then expired.
		sellerAsset1Bal := app.AssetKeeper.Balance(ctx, asset1.ID, seller)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(9990), sellerAsset1Bal)
		// 10020 because two orders executed at clearing price 2:
		// 10000 + 5 * 2 + 5 * 2 = 10020
		sellerAsset2Bal := app.AssetKeeper.Balance(ctx, asset2.ID, seller)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(10020), sellerAsset2Bal)

		buyerAsset1Bal := app.AssetKeeper.Balance(ctx, asset1.ID, buyer)
		// the orders with prices 1 and 2 receives partial fills of 5
		// the other orders expired.
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(10010), buyerAsset1Bal)

		buyerAsset2Bal := app.AssetKeeper.Balance(ctx, asset2.ID, buyer)
		// clearing of 2. two of buyer's orders were rationed for a total of 10
		// asset 2 credited.
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(9980), buyerAsset2Bal)
	})
}
