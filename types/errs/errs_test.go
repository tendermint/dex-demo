package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/dex-demo/testutil/testflags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestErrOrBlankResult(t *testing.T) {
	testflags.UnitTest(t)
	err := ErrNotFound("not found")
	assert.EqualValues(t, err.Result(), ErrOrBlankResult(err))
	assert.EqualValues(t, sdk.Result{}, ErrOrBlankResult(nil))
}
