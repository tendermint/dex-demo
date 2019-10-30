package conv

import (
	"encoding/binary"
	"io"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SDKUint2Big(in sdk.Uint) *big.Int {
	out, _ := new(big.Int).SetString(in.String(), 10)
	return out
}

func Uint642Bytes(in uint64) []byte {
	b := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, in)
	return b
}

func ReadUint64(r io.Reader) (uint64, error) {
	b := make([]byte, 8, 8)
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}
