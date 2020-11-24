package storeutil

/*
Package storeutil contains various helper methods that make working with
Cosmos SDK KVStores simpler.

Most of the methods are self-explanatory, however Increment/GetSeq()
deserve special attention. These methods increment and get an sdk.Uint
stored in the KVStore for the purposes of mimicking an auto-increment column
in SQL.

IncrementSeq takes three arguments: an sdk.Context, an  sdk.StoreKey, and
a byte slice that represents the key at which the sdk.Uint is being stored.
Upon execution, it will retrieve that Uint, increment it, and store the
incremented version in the database. GetSeq() is similar, except it just
returns the sdk.Uint.

This is useful in situations where rows in an sdk.KVStore are assigned
an automatically-incrementing sdk.Uint at creation time.
 */