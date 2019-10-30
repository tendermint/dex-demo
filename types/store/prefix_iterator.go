package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func KVStorePrefixIterator(kvs KVStore, prefix []byte) sdk.Iterator {
	return kvs.Iterator(prefix, sdk.PrefixEndBytes(prefix))
}

func KVStoreReversePrefixIterator(kvs KVStore, prefix []byte) sdk.Iterator {
	return kvs.ReverseIterator(prefix, sdk.PrefixEndBytes(prefix))
}
