package balance

import (
	"github.com/tendermint/dex-demo/embedded"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GetQueryRequest struct {
	Address sdk.AccAddress
}

type GetQueryResponseBalance struct {
	AssetID store.EntityID `json:"asset_id"`
	Name    string         `json:"name"`
	Symbol  string         `json:"symbol"`
	Liquid  sdk.Uint       `json:"liquid"`
	AtRisk  sdk.Uint       `json:"at_risk"`
}

type GetQueryResponse struct {
	Balances []GetQueryResponseBalance `json:"balances"`
}

type TransferBalanceRequest struct {
	To      sdk.AccAddress `json:"to"`
	AssetID store.EntityID `json:"asset_id"`
	Amount  sdk.Uint       `json:"amount"`
}

type TransferBalanceResponse struct {
	BlockInclusion embedded.BlockInclusion `json:"block_inclusion"`
}
