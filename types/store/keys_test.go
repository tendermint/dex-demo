package store

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/dex-demo/testutil/testflags"
)

func TestPrefixKey(t *testing.T) {
	testflags.UnitTest(t)
	out1 := PrefixKeyString("fooprefix")
	assert.Equal(t, "fooprefix", string(out1))
	out2 := PrefixKeyString("fooprefix", []byte("sub1"), []byte("sub2"))
	assert.Equal(t, "fooprefix/sub1/sub2", string(out2))
}
