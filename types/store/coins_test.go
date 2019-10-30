package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFormatCoin(t *testing.T) {
	out := FormatCoin(NewEntityID(1), sdk.NewUint(100000))
	assert.True(t, out.Amount.Equal(sdk.NewInt(100000)))
	assert.Equal(t, "asset1", out.Denom)
}

func TestFormatDenom(t *testing.T) {
	assert.Equal(t, "asset99", FormatDenom(NewEntityID(99)))
}
