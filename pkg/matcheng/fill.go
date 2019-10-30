package matcheng

import (
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Fill struct {
	OrderID     store.EntityID
	QtyFilled   sdk.Uint
	QtyUnfilled sdk.Uint
}
