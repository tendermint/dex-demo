package testutil

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AssetDecimals = 18
)

func ToBaseUnits(n uint64) sdk.Uint {
	return ToBaseUnitsDecimals(n, 0)
}

func ToBaseUnitsDecimals(n uint64, decimals int) sdk.Uint {
	return sdk.NewUint(n).Mul(sdk.NewUint(uint64(math.Pow(10, float64(AssetDecimals-decimals)))))
}
