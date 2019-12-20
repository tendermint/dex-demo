package fill

import (
	"github.com/tendermint/dex-demo/storeutil"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/embedded/store"
	"github.com/tendermint/dex-demo/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TableKey = "fill"
)

type IteratorCB func(fill Fill) bool

type Keeper struct {
	as  store.ArchiveStore
	cdc *codec.Codec
}

func NewKeeper(db dbm.DB, cdc *codec.Codec) Keeper {
	return Keeper{
		as:  store.NewTable(db, TableKey),
		cdc: cdc,
	}
}

func (k Keeper) OnFillEvent(event types.Fill) {
	fill := Fill{
		OrderID:     event.OrderID,
		Owner:       event.Owner,
		Pair:        event.Pair,
		Direction:   event.Direction,
		QtyFilled:   event.QtyFilled,
		QtyUnfilled: event.QtyUnfilled,
		BlockNumber: event.BlockNumber,
		Price:       event.Price,
	}
	storedB := k.cdc.MustMarshalBinaryBare(fill)
	k.as.Set(fillKey(event.BlockNumber, event.OrderID), storedB)
}

func (k Keeper) IterOverBlockNumbers(start int64, end int64, cb IteratorCB) {
	k.as.Iterator(fillIterKey(start), fillIterKey(end), func(_ []byte, v []byte) bool {
		var fill Fill
		k.cdc.MustUnmarshalBinaryBare(v, &fill)
		return cb(fill)
	})
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Fill:
		k.OnFillEvent(ev)
	}

	return nil
}

func fillIterKey(blockNum int64) []byte {
	return storeutil.PrefixKeyBytes(storeutil.Int64Subkey(blockNum))
}

func fillKey(blockNum int64, orderId sdk.Uint) []byte {
	return storeutil.PrefixKeyBytes(fillIterKey(blockNum), storeutil.SDKUintSubkey(orderId))
}
