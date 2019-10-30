package order

import (
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/asset"
	assettypes "github.com/tendermint/dex-demo/x/asset/types"
	"github.com/tendermint/dex-demo/x/market"
	types3 "github.com/tendermint/dex-demo/x/order/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	seqKey = "seq"
	valKey = "val"
)

type IteratorCB func(order types3.Order) bool

type Keeper struct {
	bankKeeper   bank.Keeper
	marketKeeper market.Keeper
	assetKeeper  asset.Keeper
	storeKey     sdk.StoreKey
	queue        types.Backend
	cdc          *codec.Codec
}

func NewKeeper(bk bank.Keeper, mk market.Keeper, ak asset.Keeper, sk sdk.StoreKey, queue types.Backend, cdc *codec.Codec) Keeper {
	return Keeper{
		bankKeeper:   bk,
		marketKeeper: mk,
		assetKeeper:  ak,
		storeKey:     sk,
		queue:        queue,
		cdc:          cdc,
	}
}

func (k Keeper) Post(ctx sdk.Context, owner sdk.AccAddress, mktID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16) (types3.Order, sdk.Error) {
	var err sdk.Error
	mkt, err := k.marketKeeper.Get(ctx, mktID)
	if err != nil {
		return types3.Order{}, err
	}

	var postedAsset assettypes.Asset
	var postedAmt sdk.Uint
	if direction == matcheng.Bid {
		postedAsset, err = k.assetKeeper.Get(ctx, mkt.QuoteAssetID)
		p, err := matcheng.NormalizeQuoteQuantity(price, quantity)
		if err != nil {
			return types3.Order{}, sdk.ErrInvalidCoins(err.Error())
		}
		postedAmt = p
	} else {
		postedAsset, err = k.assetKeeper.Get(ctx, mkt.BaseAssetID)
		postedAmt = quantity
	}
	if err != nil {
		// should never happen; implies consensus
		// or storage bug
		panic(err)
	}

	_, err = k.bankKeeper.SubtractCoins(ctx, owner, assettypes.Coins(postedAsset.ID, postedAmt))
	if err != nil {
		return types3.Order{}, err
	}

	return k.Create(
		ctx,
		owner,
		mktID,
		direction,
		price,
		quantity,
		tif,
	)
}

func (k Keeper) Create(ctx sdk.Context, owner sdk.AccAddress, marketID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16) (types3.Order, sdk.Error) {
	id := k.incrementSeq(ctx)
	order := types3.Order{
		ID:                id,
		Owner:             owner,
		MarketID:          marketID,
		Direction:         direction,
		Price:             price,
		Quantity:          quantity,
		TimeInForceBlocks: tif,
		CreatedBlock:      ctx.BlockHeight(),
	}
	err := store.SetNotExists(ctx, k.storeKey, k.cdc, orderKey(id), order)
	_ = k.queue.Publish(types.OrderCreated{
		ID:                order.ID,
		Owner:             order.Owner,
		MarketID:          order.MarketID,
		Direction:         order.Direction,
		Price:             order.Price,
		Quantity:          order.Quantity,
		TimeInForceBlocks: order.TimeInForceBlocks,
		CreatedBlock:      order.CreatedBlock,
	})

	return order, err
}

func (k Keeper) Cancel(ctx sdk.Context, id store.EntityID) sdk.Error {
	var err sdk.Error
	ord, err := k.Get(ctx, id)
	if err != nil {
		return err
	}
	mkt, err := k.marketKeeper.Get(ctx, ord.MarketID)
	if err != nil {
		// should never happen; implies consensus
		// or storage bug
		panic(err)
	}

	var postedAsset assettypes.Asset
	var postedAmt sdk.Uint
	if ord.Direction == matcheng.Bid {
		postedAsset, err = k.assetKeeper.Get(ctx, mkt.QuoteAssetID)
		p, err := matcheng.NormalizeQuoteQuantity(ord.Price, ord.Quantity)
		if err != nil {
			return sdk.ErrInvalidCoins(err.Error())
		}
		postedAmt = p
	} else {
		postedAsset, err = k.assetKeeper.Get(ctx, mkt.BaseAssetID)
		postedAmt = ord.Quantity
	}
	if err != nil {
		// should never happen; implies consensus
		// or storage bug
		panic(err)
	}

	_, err = k.bankKeeper.AddCoins(ctx, ord.Owner, assettypes.Coins(postedAsset.ID, postedAmt))
	if err != nil {
		// should never happen, implies consensus
		// or storage bug
		panic(err)
	}
	_ = k.queue.Publish(types.OrderCancelled{
		OrderID: id,
	})

	return k.Del(ctx, ord.ID)
}

func (k Keeper) Get(ctx sdk.Context, id store.EntityID) (types3.Order, sdk.Error) {
	var out types3.Order
	err := store.Get(ctx, k.storeKey, k.cdc, orderKey(id), &out)
	return out, err
}

func (k Keeper) Set(ctx sdk.Context, order types3.Order) sdk.Error {
	return store.SetExists(ctx, k.storeKey, k.cdc, orderKey(order.ID), order)
}

func (k Keeper) Has(ctx sdk.Context, id store.EntityID) bool {
	return store.Has(ctx, k.storeKey, orderKey(id))
}

func (k Keeper) Del(ctx sdk.Context, id store.EntityID) sdk.Error {
	return store.Del(ctx, k.storeKey, orderKey(id))
}

func (k Keeper) incrementSeq(ctx sdk.Context) store.EntityID {
	return store.IncrementSeq(ctx, k.storeKey, []byte(seqKey))
}

func (k Keeper) Iterator(ctx sdk.Context, cb IteratorCB) {
	kv := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(kv, []byte(valKey))
	k.doIterator(iter, cb)
}

func (k Keeper) ReverseIterator(ctx sdk.Context, cb IteratorCB) {
	kv := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(kv, []byte(valKey))
	k.doIterator(iter, cb)
}

func (k Keeper) doIterator(iter sdk.Iterator, cb IteratorCB) {
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		orderB := iter.Value()
		var order types3.Order
		k.cdc.MustUnmarshalBinaryBare(orderB, &order)

		if !cb(order) {
			break
		}
	}
}

func orderKey(id store.EntityID) []byte {
	return store.PrefixKeyString(valKey, id.Bytes())
}
