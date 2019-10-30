package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/dex-demo/testutil/testflags"
)

func TestHexToBytes(t *testing.T) {
	testflags.UnitTest(t)
	a, err := HexToBytes("0x0101")
	require.NoError(t, err)
	b, err := HexToBytes("0101")
	require.NoError(t, err)
	assert.EqualValues(t, a, b)
	assert.Equal(t, []byte{0x01, 0x01}, a)

	_, err = HexToBytes("foo")
	assert.Error(t, err)
}
