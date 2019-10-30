package conv

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MustAccAddressFromBech32(in string) sdk.AccAddress {
	out, err := sdk.AccAddressFromBech32(in)
	if err != nil {
		panic(err)
	}
	return out
}

func AccAddressFromECDSAPubKey(in *ecdsa.PublicKey) sdk.AccAddress {
	var tmPub secp256k1.PubKeySecp256k1
	btcPub := btcec.PublicKey(*in)
	ser := btcPub.SerializeCompressed()
	copy(tmPub[:], ser)
	return sdk.AccAddress(tmPub.Address())
}
