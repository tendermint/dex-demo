package testutil

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AssertEqualUints(t *testing.T, a sdk.Uint, b sdk.Uint, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualInts(t *testing.T, a sdk.Int, b sdk.Int, msgAndArgs ...interface{}) {
	assert.Equal(t, a.String(), b.String(), msgAndArgs...)
}

func AssertEqualHex(t *testing.T, exp string, actual []byte) {
	assert.Equal(t, exp, hex.EncodeToString(actual))
}
