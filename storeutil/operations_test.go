package storeutil

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"testing"
	"time"
)

var (
	testSKName = "test"
	testKey    = []byte("testkey")
)

func upsert(ctx sdk.Context, sk sdk.StoreKey, cdc *amino.Codec, key []byte, val interface{}) {
	s := ctx.KVStore(sk)
	s.Set(key, cdc.MustMarshalBinaryBare(val))
}

func TestSetAndGet(t *testing.T) {
	cdc := amino.NewCodec()
	ctx, sk := mockApp(t)
	upsert(ctx, sk, cdc, testKey, 12345)
	var val int
	err := Get(ctx, sk, cdc, testKey, &val)
	require.NoError(t, err)
	assert.EqualValues(t, 12345, val)
	assert.Error(t, Get(ctx, sk, cdc, []byte("not here"), &val))
}

func TestSetNotExists(t *testing.T) {
	cdc := amino.NewCodec()
	ctx, sk := mockApp(t)
	assert.NoError(t, Create(ctx, sk, cdc, testKey, 12345))
	assert.Error(t, Create(ctx, sk, cdc, testKey, 12345))
}

func TestSetExists(t *testing.T) {
	cdc := amino.NewCodec()
	ctx, sk := mockApp(t)
	assert.Error(t, Update(ctx, sk, cdc, testKey, 12345))
	upsert(ctx, sk, cdc, testKey, 12345)
	assert.NoError(t, Update(ctx, sk, cdc, testKey, 54321))
}

func TestDel(t *testing.T) {
	cdc := amino.NewCodec()
	ctx, sk := mockApp(t)
	assert.Error(t, Del(ctx, sk, testKey))
	upsert(ctx, sk, cdc, testKey, 12345)
	assert.NoError(t, Del(ctx, sk, testKey))
	assert.False(t, Has(ctx, sk, testKey))
}

func TestHas(t *testing.T) {
	cdc := amino.NewCodec()
	ctx, sk := mockApp(t)
	assert.False(t, Has(ctx, sk, testKey))
	upsert(ctx, sk, cdc, testKey, 12345)
	assert.True(t, Has(ctx, sk, testKey))
}

func mockApp(t *testing.T) (sdk.Context, sdk.StoreKey) {
	keys := sdk.NewKVStoreKeys(testSKName)
	ms := store.NewCommitMultiStore(db.NewMemDB())
	ms.MountStoreWithDB(keys[testSKName], sdk.StoreTypeIAVL, db.NewMemDB())
	require.NoError(t, ms.LoadVersion(0))
	hdr := abci.Header{ChainID: "unit-test-chain", Height: 1, Time: time.Unix(1558332092, 0)}
	return sdk.NewContext(ms, hdr, false, log.NewNopLogger()), keys[testSKName]
}

func TestIncrementSeq(t *testing.T) {
	ctx, sk := mockApp(t)
	assert.True(t, sdk.OneUint().Equal(IncrementSeq(ctx, sk, testKey)))
	assert.True(t, sdk.NewUint(2).Equal(IncrementSeq(ctx, sk, testKey)))
}

func TestGetSeq(t *testing.T) {
	ctx, sk := mockApp(t)
	assert.True(t, sdk.ZeroUint().Equal(GetSeq(ctx, sk, testKey)))
	IncrementSeq(ctx, sk, testKey)
	assert.True(t, sdk.OneUint().Equal(GetSeq(ctx, sk, testKey)))
}
