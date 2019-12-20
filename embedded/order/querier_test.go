package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/testflags"
	"github.com/tendermint/dex-demo/types"
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
		id := sdk.ZeroUint()

		for i := 0; i < 55; i++ {
			id = id.Add(sdk.OneUint())
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: sdk.NewUint(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualUints(t, sdk.NewUint(55), res.Orders[0].ID)
		testutil.AssertEqualUints(t, sdk.NewUint(6), res.Orders[49].ID)
		testutil.AssertEqualUints(t, sdk.NewUint(5), res.NextID)
	})
	t.Run("should work with an offset", func(t *testing.T) {
		id := sdk.ZeroUint()

		for i := 0; i < 55; i++ {
			id = id.Add(sdk.OneUint())
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: sdk.NewUint(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: sdk.NewUint(7),
		})
		require.NoError(t, err)

		assert.Equal(t, 7, len(res.Orders))
		testutil.AssertEqualUints(t, sdk.NewUint(7), res.Orders[0].ID)
		testutil.AssertEqualUints(t, sdk.OneUint(), res.Orders[6].ID)
		testutil.AssertEqualUints(t, sdk.ZeroUint(), res.NextID)
	})
	t.Run("should support filter by address alongside offset", func(t *testing.T) {
		id := sdk.ZeroUint()
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Add(sdk.OneUint())
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: sdk.NewUint(2),
				ID:       id,
				Owner:    owner,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: sdk.NewUint(104),
			Owner: genOwner,
		})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualUints(t, sdk.NewUint(109), res.Orders[0].ID)
		testutil.AssertEqualUints(t, sdk.NewUint(11), res.Orders[49].ID)
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
