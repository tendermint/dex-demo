package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestQuerier(t *testing.T) {
	testflags.UnitTest(t)
	cdc := codec.New()
	db := dbm.NewMemDB()
	k := NewKeeper(db, cdc)
	ctx := testutil.DummyContext()
	q := NewQuerier(k)
	doListQuery := func(req ListQueryRequest) (ListQueryResult, error) {
		reqB := serializeRequestQuery(cdc, req)
		var res ListQueryResult
		resB, err := q(ctx, []string{"list"}, reqB)
		if err != nil {
			return res, err
		}
		cdc.MustUnmarshalJSON(resB, &res)
		return res, nil
	}

	t.Run("should return no more than 50 orders in descending order", func(t *testing.T) {
		id := store.NewEntityID(0)

		for i := 0; i < 55; i++ {
			id = id.Inc()
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(55), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(6), res.Orders[49].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(5), res.NextID)
	})
	t.Run("should work with an offset", func(t *testing.T) {
		id := store.NewEntityID(0)

		for i := 0; i < 55; i++ {
			id = id.Inc()
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: store.NewEntityID(7),
		})
		require.NoError(t, err)

		assert.Equal(t, 7, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(7), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(1), res.Orders[6].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(0), res.NextID)
	})
	t.Run("should support filter by address alongside offset", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
				Owner:    owner,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: store.NewEntityID(104),
			Owner: genOwner,
		})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(109), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(11), res.Orders[49].ID)
	})
	t.Run("should return an error if the request does not deserialize", func(t *testing.T) {
		_, err := q(ctx, []string{"list"}, abci.RequestQuery{Data: []byte("foo")})
		require.Error(t, err)
	})
}

func serializeRequestQuery(cdc *codec.Codec, req ListQueryRequest) abci.RequestQuery {
	data := cdc.MustMarshalBinaryBare(req)

	return abci.RequestQuery{
		Data: data,
	}
}
