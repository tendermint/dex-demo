package matcheng

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Fill struct {
	OrderID     sdk.Uint
	QtyFilled   sdk.Uint
	QtyUnfilled sdk.Uint
}
