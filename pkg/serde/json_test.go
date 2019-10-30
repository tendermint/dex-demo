package serde

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/dex-demo/testutil/testflags"
)

type example struct {
	Value HexBytes `json:"value"`
}

func TestHexBytes(t *testing.T) {
	testflags.UnitTest(t)
	ser := example{Value: []byte{0x99}}
	out, err := json.Marshal(ser)
	require.NoError(t, err)
	assert.Equal(t, "{\"value\":\"0x99\"}", string(out))
	var deser example
	err = json.Unmarshal(out, &deser)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(ser.Value, deser.Value))

	ser = example{Value: nil}
	out, err = json.Marshal(ser)
	require.NoError(t, err)
	assert.Equal(t, "{\"value\":null}", string(out))
	err = json.Unmarshal(out, &deser)
	require.NoError(t, err)
	assert.Nil(t, deser.Value)

	err = json.Unmarshal([]byte("{\"value\":\"}"), &deser)
	assert.Error(t, err)
	err = json.Unmarshal([]byte("{\"value\":\"\"}"), &deser)
	assert.Error(t, err)
}
