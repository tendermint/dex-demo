package testutil

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AssertEqualUints(t *testing.T, a sdk.Uint, b sdk.Uint, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualInts(t *testing.T, a sdk.Int, b sdk.Int, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualEntityIDs(t *testing.T, a store.EntityID, b store.EntityID, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualHex(t *testing.T, exp string, actual []byte) {
	assert.Equal(t, exp, hex.EncodeToString(actual))
}
