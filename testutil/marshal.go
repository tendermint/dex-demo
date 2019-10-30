package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func MustUnmarshalJSON(t *testing.T, data []byte, proto interface{}) {
	err := amino.UnmarshalJSON(data, proto)
	require.NoError(t, err)
}
