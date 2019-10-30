package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper(t *testing.T) {
	testflags.UnitTest(t)
	cdc := codec.New()
	db := dbm.NewMemDB()
	k := NewKeeper(db, cdc)
	creationEvs := []types.OrderCreated{
		{
			ID:                store.NewEntityID(1),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(1),
			Direction:         matcheng.Bid,
			Price:             sdk.NewUint(100),
			Quantity:          sdk.NewUint(100),
			TimeInForceBlocks: 19,
			CreatedBlock:      10,
		},
		{
			ID:                store.NewEntityID(2),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(1),
			Direction:         matcheng.Ask,
			Price:             sdk.NewUint(110),
			Quantity:          sdk.NewUint(110),
			TimeInForceBlocks: 20,
			CreatedBlock:      11,
		},
		{
			ID:                store.NewEntityID(3),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(2),
			Direction:         matcheng.Bid,
			Price:             sdk.NewUint(99),
			Quantity:          sdk.NewUint(99),
			TimeInForceBlocks: 20,
			CreatedBlock:      12,
		},
		{
			ID:                store.NewEntityID(4),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(1),
			Direction:         matcheng.Bid,
			Price:             sdk.NewUint(100),
			Quantity:          sdk.NewUint(100),
			TimeInForceBlocks: 19,
			CreatedBlock:      10,
		},
		{
			ID:                store.NewEntityID(5),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(1),
			Direction:         matcheng.Bid,
			Price:             sdk.NewUint(100),
			Quantity:          sdk.NewUint(100),
			TimeInForceBlocks: 19,
			CreatedBlock:      10,
		},
		{
			ID:                store.NewEntityID(6),
			Owner:             sdk.AccAddress{},
			MarketID:          store.NewEntityID(2),
			Direction:         matcheng.Bid,
			Price:             sdk.NewUint(100),
			Quantity:          sdk.NewUint(100),
			TimeInForceBlocks: 19,
			CreatedBlock:      10,
		},
	}
	cancellationEvs := []types.OrderCancelled{
		{
			OrderID: store.NewEntityID(4),
		},
	}
	fillEvs := []types.Fill{
		{
			OrderID:   store.NewEntityID(5),
			QtyFilled: sdk.NewUint(99),
		},
		{
			OrderID:   store.NewEntityID(6),
			QtyFilled: sdk.NewUint(100),
		},
	}
	for _, e := range creationEvs {
		require.NoError(t, k.OnEvent(e))
	}
	for _, e := range cancellationEvs {
		require.NoError(t, k.OnEvent(e))
	}
	for _, e := range fillEvs {
		require.NoError(t, k.OnEvent(e))
	}

	t.Run("open orders by market returns only open orders from the market", func(t *testing.T) {
		res := k.OpenOrdersByMarket(store.NewEntityID(1))
		assert.Equal(t, 3, len(res))
		ev0 := creationEvs[0]
		ev1 := creationEvs[1]
		ev4 := creationEvs[4]
		assertEqualOrders(t, cdc, Order{
			ID:             ev4.ID,
			Owner:          ev4.Owner,
			MarketID:       ev4.MarketID,
			Direction:      ev4.Direction,
			Price:          ev4.Price,
			Quantity:       ev4.Quantity,
			Status:         "OPEN",
			Type:           "LIMIT",
			TimeInForce:    ev4.TimeInForceBlocks,
			QuantityFilled: sdk.NewUint(99),
			CreatedBlock:   ev4.CreatedBlock,
		}, res[0])
		assertEqualOrders(t, cdc, Order{
			ID:             ev1.ID,
			Owner:          ev1.Owner,
			MarketID:       ev1.MarketID,
			Direction:      ev1.Direction,
			Price:          ev1.Price,
			Quantity:       ev1.Quantity,
			Status:         "OPEN",
			Type:           "LIMIT",
			TimeInForce:    ev1.TimeInForceBlocks,
			QuantityFilled: sdk.NewUint(0),
			CreatedBlock:   ev1.CreatedBlock,
		}, res[1])
		assertEqualOrders(t, cdc, Order{
			ID:             ev0.ID,
			Owner:          ev0.Owner,
			MarketID:       ev0.MarketID,
			Direction:      ev0.Direction,
			Price:          ev0.Price,
			Quantity:       ev0.Quantity,
			Status:         "OPEN",
			Type:           "LIMIT",
			TimeInForce:    ev0.TimeInForceBlocks,
			QuantityFilled: sdk.NewUint(0),
			CreatedBlock:   ev0.CreatedBlock,
		}, res[2])
	})
	t.Run("cancelled orders are returned as cancelled", func(t *testing.T) {
		ev3 := creationEvs[3]
		res, err := k.Get(ev3.ID)
		require.NoError(t, err)
		assertEqualOrders(t, cdc, Order{
			ID:             ev3.ID,
			Owner:          ev3.Owner,
			MarketID:       ev3.MarketID,
			Direction:      ev3.Direction,
			Price:          ev3.Price,
			Quantity:       ev3.Quantity,
			Status:         "CANCELLED",
			Type:           "LIMIT",
			TimeInForce:    ev3.TimeInForceBlocks,
			QuantityFilled: sdk.NewUint(0),
			CreatedBlock:   ev3.CreatedBlock,
		}, res)
	})
	t.Run("fully filled orders are returned as filled", func(t *testing.T) {
		ev5 := creationEvs[5]
		res, err := k.Get(ev5.ID)
		require.NoError(t, err)
		assertEqualOrders(t, cdc, Order{
			ID:             ev5.ID,
			Owner:          ev5.Owner,
			MarketID:       ev5.MarketID,
			Direction:      ev5.Direction,
			Price:          ev5.Price,
			Quantity:       ev5.Quantity,
			Status:         "FILLED",
			Type:           "LIMIT",
			TimeInForce:    ev5.TimeInForceBlocks,
			QuantityFilled: sdk.NewUint(100),
			CreatedBlock:   ev5.CreatedBlock,
		}, res)
	})
}

func assertEqualOrders(t *testing.T, cdc *codec.Codec, exp Order, actual Order) {
	expJson := cdc.MustMarshalJSON(exp)
	actualJson := cdc.MustMarshalJSON(actual)
	assert.Equal(t, string(expJson), string(actualJson))
}
