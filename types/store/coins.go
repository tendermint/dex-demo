package store

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FormatCoin(id EntityID, amount sdk.Uint) sdk.Coin {
	out, err := sdk.ParseCoin(fmt.Sprintf("%s%s", amount.String(), FormatDenom(id)))
	// should never happen
	if err != nil {
		panic(err)
	}
	return out
}

func FormatDenom(id EntityID) string {
	return fmt.Sprintf("asset%s", id.String())
}
