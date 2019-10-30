package types

import (
	"github.com/tendermint/dex-demo/types/store"
)

const (
	ModuleName = "market"
	RouterKey  = ModuleName
	StoreKey   = RouterKey
)

type Market struct {
	ID           store.EntityID
	BaseAssetID  store.EntityID
	QuoteAssetID store.EntityID
}

func New(id store.EntityID, baseAsset store.EntityID, quoteAsset store.EntityID) Market {
	return Market{
		ID:           id,
		BaseAssetID:  baseAsset,
		QuoteAssetID: quoteAsset,
	}
}
