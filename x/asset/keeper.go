package asset

import (
	"github.com/tendermint/dex-demo/embedded/store"
	"github.com/tendermint/dex-demo/pkg/log"
	"github.com/tendermint/dex-demo/storeutil"
	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/x/asset/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	seqKey = "seq"
	valKey = "val"
)

type IteratorCB func(asset types.Asset) bool

type Keeper struct {
	bankKeeper bank.Keeper

	storeKey sdk.StoreKey

	cdc *codec.Codec
}

var logger = log.WithModule("order_keeper")

func NewKeeper(bk bank.Keeper, sk sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		bankKeeper: bk,
		storeKey:   sk,
		cdc:        cdc,
	}
}

func (k Keeper) Create(ctx sdk.Context, name string, symbol string, owner sdk.AccAddress, totalSupply sdk.Uint) (types.Asset, sdk.Error) {
	id := k.incrementSeq(ctx)
	asset := types.New(
		id,
		name,
		symbol,
		owner,
		sdk.ZeroUint(),
		totalSupply,
	)
	err := errs.WrapOrNil(storeutil.Create(ctx, k.storeKey, k.cdc, assetKey(id), asset))
	return asset, err
}

func (k Keeper) Inject(ctx sdk.Context, asset types.Asset) {
	seq := storeutil.GetSeq(ctx, k.storeKey, []byte(seqKey))

	if !asset.ID.Sub(sdk.OneUint()).Equal(seq) {
		panic("Invalid asset ID.")
	}

	k.incrementSeq(ctx)
	if err := storeutil.Create(ctx, k.storeKey, k.cdc, assetKey(asset.ID), asset); err != nil {
		panic(err)
	}
}

func (k Keeper) Set(ctx sdk.Context, asset types.Asset) sdk.Error {
	return errs.WrapNotFound(storeutil.Update(ctx, k.storeKey, k.cdc, assetKey(asset.ID), asset))
}

func (k Keeper) Get(ctx sdk.Context, id sdk.Uint) (types.Asset, sdk.Error) {
	var a types.Asset
	err := errs.WrapNotFound(storeutil.Get(ctx, k.storeKey, k.cdc, assetKey(id), &a))
	return a, err
}

func (k Keeper) Has(ctx sdk.Context, id sdk.Uint) bool {
	return storeutil.Has(ctx, k.storeKey, assetKey(id))
}

func (k Keeper) Iterator(ctx sdk.Context, cb IteratorCB) {
	kv := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(kv, []byte(valKey))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		aB := iter.Value()
		var asset types.Asset
		k.cdc.MustUnmarshalBinaryBare(aB, &asset)

		if !cb(asset) {
			break
		}
	}
}

func (k Keeper) Mint(ctx sdk.Context, id sdk.Uint, quantity sdk.Uint) sdk.Error {
	asset, err := k.Get(ctx, id)
	if err != nil {
		return err
	}
	newSupply := asset.CirculatingSupply.Add(quantity)
	if newSupply.GT(asset.TotalSupply) {
		logger.Info(
			"rejected mint for more than total supply",
			"new_supply",
			newSupply.String(),
			"circulating_supply",
			asset.CirculatingSupply.String(),
			"quantity",
			quantity.String(),
		)
		return sdk.ErrInvalidCoins("cannot mint more than total supply")
	}
	_, err = k.bankKeeper.AddCoins(ctx, asset.Owner, types.Coins(asset.ID, quantity))
	if err != nil {
		return err
	}

	asset.CirculatingSupply = newSupply
	return k.Set(ctx, asset)
}

func (k Keeper) Burn(ctx sdk.Context, id sdk.Uint, burner sdk.AccAddress, quantity sdk.Uint) sdk.Error {
	asset, err := k.Get(ctx, id)
	if err != nil {
		return err
	}
	if asset.CirculatingSupply.LT(quantity) {
		return sdk.ErrInvalidCoins("cannot burn more than circulating supply")
	}
	newSupply := asset.CirculatingSupply.Sub(quantity)
	_, err = k.bankKeeper.SubtractCoins(ctx, burner, types.Coins(asset.ID, quantity))
	if err != nil {
		return err
	}

	asset.CirculatingSupply = newSupply
	return k.Set(ctx, asset)
}

func (k Keeper) Transfer(ctx sdk.Context, id sdk.Uint, from sdk.AccAddress, to sdk.AccAddress, quantity sdk.Uint) sdk.Error {
	asset, err := k.Get(ctx, id)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoins(ctx, from, to, types.Coins(asset.ID, quantity))
}

func (k Keeper) TransferFromOwner(ctx sdk.Context, id sdk.Uint, to sdk.AccAddress, quantity sdk.Uint) sdk.Error {
	asset, err := k.Get(ctx, id)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoins(ctx, asset.Owner, to, types.Coins(asset.ID, quantity))
}

func (k Keeper) Balance(ctx sdk.Context, id sdk.Uint, owner sdk.AccAddress) sdk.Uint {
	coins := k.bankKeeper.GetCoins(ctx, owner)
	out := coins.AmountOf(store.FormatDenom(id))
	return sdk.NewUintFromBigInt(out.BigInt())
}

func (k Keeper) incrementSeq(ctx sdk.Context) sdk.Uint {
	assetNum := storeutil.IncrementSeq(ctx, k.storeKey, []byte(seqKey))
	return assetNum
}

func assetKey(id sdk.Uint) []byte {
	return storeutil.PrefixKeyString(valKey, storeutil.SDKUintSubkey(id))
}
