package storeutil

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ErrStoreKeyNotFound = errors.New("key not found")
	ErrKeyExists        = errors.New("key exists")
)

// Get unmarshals a binary object in the store
// identified by sk and key into the object
// identified by proto.
func Get(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, proto interface{}) error {
	store := ctx.KVStore(sk)
	b := store.Get(key)
	if b == nil {
		return ErrStoreKeyNotFound
	}
	cdc.MustUnmarshalBinaryBare(b, proto)
	return nil
}

// Create inserts val into the store
// identified by sk and at the key
// identified by key. Create will return
// an error if the key already exists.
func Create(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, val interface{}) error {
	if Has(ctx, sk, key) {
		return ErrKeyExists
	}
	store := ctx.KVStore(sk)
	store.Set(key, cdc.MustMarshalBinaryBare(val))
	return nil
}

// Update inserts val into the store
// identified by sk and at the key
// identified by key. Update will return
// an error if the key does not exist.
func Update(ctx sdk.Context, sk sdk.StoreKey, cdc *codec.Codec, key []byte, val interface{}) error {
	if !Has(ctx, sk, key) {
		return ErrStoreKeyNotFound
	}
	store := ctx.KVStore(sk)
	store.Set(key, cdc.MustMarshalBinaryBare(val))
	return nil
}

// Del deletes the value in the store
// identified by sk and at the key
// identified by key. Del will return an error
// if the key does not exist.
func Del(ctx sdk.Context, sk sdk.StoreKey, key []byte) error {
	if !Has(ctx, sk, key) {
		return ErrStoreKeyNotFound
	}
	store := ctx.KVStore(sk)
	store.Delete(key)
	return nil
}

// Has returns true if the specified key
// exists in the store identified by sk.
func Has(ctx sdk.Context, sk sdk.StoreKey, key []byte) bool {
	store := ctx.KVStore(sk)
	return store.Has(key)
}

// IncrementSeq increments the Uint in the store
// identified by sk at the key seqKey.
func IncrementSeq(ctx sdk.Context, sk sdk.StoreKey, seqKey []byte) sdk.Uint {
	store := ctx.KVStore(sk)
	seq := GetSeq(ctx, sk, seqKey).Add(sdk.OneUint())
	store.Set(seqKey, []byte(seq.String()))
	return seq
}

// GetSeq returns the Uint in the store
// identified by sk at the key seqKey.
func GetSeq(ctx sdk.Context, sk sdk.StoreKey, seqKey []byte) sdk.Uint {
	store := ctx.KVStore(sk)
	if !store.Has(seqKey) {
		return sdk.ZeroUint()
	}

	b := store.Get(seqKey)
	return sdk.NewUintFromString(string(b))
}