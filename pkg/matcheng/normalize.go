package matcheng

import (
	"errors"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AssetDecimals = 18
)

var divisor = sdk.NewDec(int64(math.Pow(float64(10), float64(AssetDecimals))))

func NormalizeQuoteQuantity(quotePrice sdk.Uint, baseQuantity sdk.Uint) (sdk.Uint, error) {
	quotePDec := sdk.NewDecFromBigInt(quotePrice.BigInt())
	baseQDec := sdk.NewDecFromBigInt(baseQuantity.BigInt())
	baseMult := baseQDec.Quo(divisor)
	res := sdk.NewUintFromBigInt(quotePDec.Mul(baseMult).TruncateInt().BigInt())
	var err error
	if res.IsZero() {
		err = errors.New("quantity too small to represent")
	}
	return res, err
}
