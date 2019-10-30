package order_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/mockapp"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/asset/types"
	types2 "github.com/tendermint/dex-demo/x/market/types"
	types4 "github.com/tendermint/dex-demo/x/order/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type testCtx struct {
	ctx      sdk.Context
	marketID store.EntityID
	owner    sdk.AccAddress
	buyer    sdk.AccAddress
	seller   sdk.AccAddress
	app      *mockapp.MockApp
	asset1   types.Asset
	asset2   types.Asset
	market   types2.Market
}

func TestKeeper_Post(t *testing.T) {
	testflags.UnitTest(t)
	t.Run("returns an error for a nonexistent market", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID.Inc(), matcheng.Bid, testutil.ToBaseUnits(1), testutil.ToBaseUnits(10), 599)
		assert.Error(t, err)
		assert.Equal(t, err.Code(), errs.CodeNotFound)
	})
	t.Run("returns an error if buying more than owned coins", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(5001), 599)
		assert.Error(t, err)
		assert.Equal(t, err.Code(), sdk.CodeInsufficientCoins)
	})
	t.Run("returns an error if selling more than owned coins", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.seller, ctx.marketID, matcheng.Ask, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10001), 599)
		assert.Error(t, err)
		assert.Equal(t, err.Code(), sdk.CodeInsufficientCoins)
	})
	t.Run("returns an error if trying to post a non-representable order", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.seller, ctx.marketID, matcheng.Bid, sdk.NewUint(2), sdk.NewUint(2), 599)
		assert.Error(t, err)
		assert.Equal(t, err.Code(), sdk.CodeInvalidCoins)
	})
	t.Run("creates the order", func(t *testing.T) {
		ctx := setupTest(t)
		created, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(1), testutil.ToBaseUnits(10), 599)
		require.NoError(t, err)
		retrieved, err := ctx.app.OrderKeeper.Get(ctx.ctx, created.ID)
		require.NoError(t, err)
		assert.EqualValues(t, created, retrieved)
	})
	t.Run("debits the correct coins", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
		require.NoError(t, err)
		_, err = ctx.app.OrderKeeper.Post(ctx.ctx, ctx.seller, ctx.marketID, matcheng.Ask, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(9980), ctx.app.AssetKeeper.Balance(ctx.ctx, ctx.market.QuoteAssetID, ctx.buyer))
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(9990), ctx.app.AssetKeeper.Balance(ctx.ctx, ctx.market.BaseAssetID, ctx.seller))
	})
}

func TestKeeper_Cancel(t *testing.T) {
	testflags.UnitTest(t)
	t.Run("returns an error for a nonexistent order", func(t *testing.T) {
		ctx := setupTest(t)
		_, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, store.NewEntityID(0), matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
		assert.Error(t, err)
		assert.Equal(t, err.Code(), errs.CodeNotFound)
	})
	t.Run("deletes the order and returns coins after cancellation", func(t *testing.T) {
		ctx := setupTest(t)
		bid, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
		require.NoError(t, err)
		err = ctx.app.OrderKeeper.Cancel(ctx.ctx, bid.ID)
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(10000), ctx.app.AssetKeeper.Balance(ctx.ctx, ctx.market.QuoteAssetID, ctx.buyer))
		assert.False(t, ctx.app.OrderKeeper.Has(ctx.ctx, bid.ID))
	})
}

func TestKeeper_Iteration(t *testing.T) {
	testflags.UnitTest(t)
	ctx := setupTest(t)
	first, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
	require.NoError(t, err)
	_, err = ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
	require.NoError(t, err)
	last, err := ctx.app.OrderKeeper.Post(ctx.ctx, ctx.buyer, ctx.marketID, matcheng.Bid, testutil.ToBaseUnits(2), testutil.ToBaseUnits(10), 599)
	require.NoError(t, err)

	var coll []store.EntityID
	ctx.app.OrderKeeper.Iterator(ctx.ctx, func(order types4.Order) bool {
		if order.ID.Equals(last.ID) {
			return false
		}
		coll = append(coll, order.ID)
		return true
	})
	assert.EqualValues(t, []store.EntityID{store.NewEntityID(1), store.NewEntityID(2)}, coll)

	coll = make([]store.EntityID, 0)
	ctx.app.OrderKeeper.ReverseIterator(ctx.ctx, func(order types4.Order) bool {
		if order.ID.Equals(first.ID) {
			return false
		}
		coll = append(coll, order.ID)
		return true
	})
	assert.EqualValues(t, []store.EntityID{store.NewEntityID(3), store.NewEntityID(2)}, coll)
}

func setupTest(t *testing.T) *testCtx {
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

	return &testCtx{
		ctx:      app.Ctx,
		marketID: mkt.ID,
		owner:    owner,
		buyer:    buyer,
		seller:   seller,
		app:      app,
		asset1:   asset1,
		asset2:   asset2,
		market:   mkt,
	}
}
