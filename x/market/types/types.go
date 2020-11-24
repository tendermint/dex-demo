package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "market"
	RouterKey  = ModuleName
	StoreKey   = RouterKey
)

type Market struct {
	ID           sdk.Uint
	BaseAssetID  sdk.Uint
	QuoteAssetID sdk.Uint
}

func New(id sdk.Uint, baseAsset sdk.Uint, quoteAsset sdk.Uint) Market {
	return Market{
		ID:           id,
		BaseAssetID:  baseAsset,
		QuoteAssetID: quoteAsset,
	}
}
