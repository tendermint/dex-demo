package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type KVStore interface {
	Get(key []byte) []byte

	Has(key []byte) bool

	Set(key, value []byte)

	Delete(key []byte)

	Iterator(start, end []byte) sdk.Iterator

	ReverseIterator(start, end []byte) sdk.Iterator
}
