package price

import (
	"github.com/tendermint/dex-demo/storeutil"
	"time"

	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/embedded/store"
	"github.com/tendermint/dex-demo/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IteratorCB func(tick Tick) bool

type Keeper struct {
	as  store.ArchiveStore
	cdc *codec.Codec
}

func NewKeeper(db dbm.DB, cdc *codec.Codec) Keeper {
	return Keeper{
		as:  store.NewTable(db, EntityName),
		cdc: cdc,
	}
}

func (k Keeper) ReverseIteratorByMarket(mktID sdk.Uint, cb IteratorCB) {
	k.as.PrefixIterator(tickIterKey(mktID), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) ReverseIteratorByMarketFrom(mktID sdk.Uint, from time.Time, cb IteratorCB) {
	k.as.ReverseIterator(tickKey(mktID, 0), sdk.PrefixEndBytes(tickKey(mktID, 0)), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) IteratorByMarketAndInterval(mktID sdk.Uint, from time.Time, to time.Time, cb IteratorCB) {
	k.as.Iterator(tickKey(mktID, from.Unix()), sdk.PrefixEndBytes(tickKey(mktID, to.Unix())), func(_ []byte, v []byte) bool {
		var tick Tick
		k.cdc.MustUnmarshalBinaryBare(v, &tick)
		return cb(tick)
	})
}

func (k Keeper) OnFillEvent(event types.Fill) {
	tick := Tick{
		MarketID:    event.MarketID,
		Pair:        event.Pair,
		BlockNumber: event.BlockNumber,
		BlockTime:   event.BlockTime,
		Price:       event.Price,
	}
	storedB := k.cdc.MustMarshalBinaryBare(tick)
	k.as.Set(tickKey(event.MarketID, tick.BlockTime), storedB)
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Fill:
		k.OnFillEvent(ev)
	}

	return nil
}

func tickKey(mktID sdk.Uint, blockTime int64) []byte {
	return storeutil.PrefixKeyBytes(tickIterKey(mktID), storeutil.Int64Subkey(blockTime))
}

func tickIterKey(mktID sdk.Uint) []byte {
	return storeutil.PrefixKeyString("tick", storeutil.SDKUintSubkey(mktID))
}
