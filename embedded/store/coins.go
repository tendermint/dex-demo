package store

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FormatCoin(id sdk.Uint, amount sdk.Uint) sdk.Coin {
	out, err := sdk.ParseCoin(fmt.Sprintf("%s%s", amount.String(), FormatDenom(id)))
	// should never happen
	if err != nil {
		panic(err)
	}
	return out
}

func FormatDenom(id sdk.Uint) string {
	return fmt.Sprintf("asset%s", id.String())
}
