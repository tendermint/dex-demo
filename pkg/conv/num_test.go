package conv

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/dex-demo/testutil/testflags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestUint2Big(t *testing.T) {
	testflags.UnitTest(t)
	a := sdk.NewUint(1)
	b := big.NewInt(1)
	assert.Equal(t, "1", SDKUint2Big(a).String())
	assert.EqualValues(t, b, SDKUint2Big(a))
}
