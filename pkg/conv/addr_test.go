package conv

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAccAddressFromECDSAPubKey(t *testing.T) {
	pub := "cosmospub1addwnpepq23xtezz0cm48uxmwlr7wkz8mh39hytc5pxzwgjahvnes8yvywh924qf238"
	decPub := sdk.MustGetAccPubKeyBech32(pub).(secp256k1.PubKeySecp256k1)
	btcPub, err := btcec.ParsePubKey(decPub[:], btcec.S256())
	require.NoError(t, err)
	addr := AccAddressFromECDSAPubKey(btcPub.ToECDSA())
	assert.Equal(t, "cosmos14hvaduk4ghcre8h44rrg4nmhx3j6wa24kj00kg", addr.String())
}
