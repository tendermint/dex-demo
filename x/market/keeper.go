package market

import (
	"fmt"

	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/asset"
	"github.com/tendermint/dex-demo/x/market/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	seqKey = "seq"
	valKey = "val"
)

type IteratorCB func(mkt types.Market) bool

type Keeper struct {
	storeKey sdk.StoreKey
	ak       asset.Keeper
	cdc      *codec.Codec
}

func NewKeeper(sk sdk.StoreKey, ak asset.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: sk,
		ak:       ak,
		cdc:      cdc,
	}
}

func (k Keeper) Create(ctx sdk.Context, baseAsset store.EntityID, quoteAsset store.EntityID) types.Market {
	id := store.IncrementSeq(ctx, k.storeKey, []byte(seqKey))
	market := types.Market{
		ID:           id,
		BaseAssetID:  baseAsset,
		QuoteAssetID: quoteAsset,
	}
	err := store.SetNotExists(ctx, k.storeKey, k.cdc, marketKey(id), market)
	// should never happen, implies consensus
	// or storage bug
	if err != nil {
		panic(err)
	}
	return market
}

func (k Keeper) Inject(ctx sdk.Context, market types.Market) {
	seq := store.GetSeq(ctx, k.storeKey, []byte(seqKey))

	if !market.ID.Dec().Equals(seq) {
		panic("Invalid asset ID.")
	}

	store.IncrementSeq(ctx, k.storeKey, []byte(seqKey))
	if err := store.SetNotExists(ctx, k.storeKey, k.cdc, marketKey(market.ID), market); err != nil {
		panic(err)
	}
}

func (k Keeper) Get(ctx sdk.Context, id store.EntityID) (types.Market, sdk.Error) {
	var m types.Market
	err := store.Get(ctx, k.storeKey, k.cdc, marketKey(id), &m)
	return m, err
}

func (k Keeper) Pair(ctx sdk.Context, id store.EntityID) (string, sdk.Error) {
	mkt, err := k.Get(ctx, id)
	if err != nil {
		return "", err
	}
	base, err := k.ak.Get(ctx, mkt.BaseAssetID)
	if err != nil {
		panic(err)
	}
	quote, err := k.ak.Get(ctx, mkt.QuoteAssetID)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/%s", base.Symbol, quote.Symbol), nil
}

func (k Keeper) Has(ctx sdk.Context, id store.EntityID) bool {
	return store.Has(ctx, k.storeKey, marketKey(id))
}

func (k Keeper) Iterator(ctx sdk.Context, cb IteratorCB) {
	kv := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(kv, []byte(valKey))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		mktB := iter.Value()
		var mkt types.Market
		k.cdc.MustUnmarshalBinaryBare(mktB, &mkt)

		if !cb(mkt) {
			break
		}
	}
}

func marketKey(id store.EntityID) []byte {
	return store.PrefixKeyString(valKey, id.Bytes())
}

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized market message type: %T", msg)).Result()
	}
}
