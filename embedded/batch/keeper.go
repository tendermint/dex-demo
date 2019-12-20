package batch

import (
	"github.com/tendermint/dex-demo/storeutil"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/embedded/store"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/errs"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TableKey = "batch"

	batchKeyPrefix = "batch"
)

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

func (k Keeper) LatestByMarket(marketID sdk.Uint) (Batch, sdk.Error) {
	var res Batch
	var found bool
	k.as.ReversePrefixIterator(batchIterKey(marketID), func(_ []byte, v []byte) bool {
		k.cdc.MustUnmarshalBinaryBare(v, &res)
		found = true
		return false
	})

	if !found {
		return res, errs.ErrNotFound("batch not found")
	}

	return res, nil
}

func (k Keeper) OnBatchEvent(event types.Batch) {
	batch := Batch{
		BlockNumber:   event.BlockNumber,
		BlockTime:     event.BlockTime,
		MarketID:      event.MarketID,
		ClearingPrice: event.ClearingPrice,
		Bids:          event.Bids,
		Asks:          event.Asks,
	}
	k.as.Set(batchKey(batch.MarketID, batch.BlockNumber), k.cdc.MustMarshalBinaryBare(batch))
}

func (k Keeper) OnEvent(event interface{}) error {
	switch ev := event.(type) {
	case types.Batch:
		k.OnBatchEvent(ev)
	}

	return nil
}

func batchKey(marketID sdk.Uint, blkNum int64) []byte {
	return storeutil.PrefixKeyBytes(batchIterKey(marketID), storeutil.Int64Subkey(blkNum))
}

func batchIterKey(marketID sdk.Uint) []byte {
	return storeutil.PrefixKeyString(batchKeyPrefix, storeutil.SDKUintSubkey(marketID))
}
