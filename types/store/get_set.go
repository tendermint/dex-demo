package store

import (
	"github.com/tendermint/dex-demo/types/errs"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Get(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, proto interface{}) sdk.Error {
	store := ctx.KVStore(sk)
	b := store.Get(key)
	if b == nil {
		return errs.ErrNotFound("not found")
	}
	cdc.MustUnmarshalBinaryBare(b, proto)
	return nil
}

func Set(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, val interface{}) {
	store := ctx.KVStore(sk)
	store.Set(key, cdc.MustMarshalBinaryBare(val))
}

func SetNotExists(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, val interface{}) sdk.Error {
	if Has(ctx, sk, key) {
		return errs.ErrAlreadyExists("already exists")
	}
	Set(ctx, sk, cdc, key, val)
	return nil
}

func SetExists(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, val interface{}) sdk.Error {
	if !Has(ctx, sk, key) {
		return errs.ErrNotFound("not found")
	}
	Set(ctx, sk, cdc, key, val)
	return nil
}

func Del(ctx sdk.Context, sk sdk.StoreKey, key []byte) sdk.Error {
	if !Has(ctx, sk, key) {
		return errs.ErrNotFound("not found")
	}
	store := ctx.KVStore(sk)
	store.Delete(key)
	return nil
}

func Has(ctx sdk.Context, sk sdk.StoreKey, key []byte) bool {
	store := ctx.KVStore(sk)
	return store.Has(key)
}

func IncrementSeq(ctx sdk.Context, sk sdk.StoreKey, seqKey []byte) EntityID {
	store := ctx.KVStore(sk)
	seq := GetSeq(ctx, sk, seqKey).Inc()
	store.Set(seqKey, []byte(seq.String()))
	return seq
}

func GetSeq(ctx sdk.Context, sk sdk.StoreKey, seqKey []byte) EntityID {
	store := ctx.KVStore(sk)
	if !store.Has(seqKey) {
		return ZeroEntityID
	}

	b := store.Get(seqKey)
	return NewEntityIDFromString(string(b))
}
