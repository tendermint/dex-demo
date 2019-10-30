package types

import (
	"fmt"
	"strings"

	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "asset"
	RouterKey  = ModuleName
	StoreKey   = RouterKey
)

type Asset struct {
	ID                store.EntityID `json:"id"`
	Name              string         `json:"name"`
	Symbol            string         `json:"symbol"`
	Owner             sdk.AccAddress `json:"owner"`
	CirculatingSupply sdk.Uint       `json:"circulating_supply"`
	TotalSupply       sdk.Uint       `json:"total_supply"`
}

func (a Asset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s",
Name: %s,
Symbol: %s,
Owner: %s,
Circulating Supply: %s,
TotalSupply: %s"`, a.ID, a.Name, a.Symbol, a.Owner, a.CirculatingSupply, a.TotalSupply))
}

func New(id store.EntityID, name string, symbol string, owner sdk.AccAddress, circSup sdk.Uint, totalSup sdk.Uint) Asset {
	return Asset{
		ID:                id,
		Name:              name,
		Symbol:            symbol,
		Owner:             owner,
		CirculatingSupply: circSup,
		TotalSupply:       totalSup,
	}
}

func Coin(id store.EntityID, quantity sdk.Uint) sdk.Coin {
	return store.FormatCoin(id, quantity)
}

func Coins(id store.EntityID, quantity sdk.Uint) sdk.Coins {
	return sdk.NewCoins(Coin(id, quantity))
}
