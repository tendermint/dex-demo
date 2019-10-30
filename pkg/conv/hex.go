package conv

import (
	"encoding/hex"
	"strings"
)

func HexToBytes(in string) ([]byte, error) {
	if strings.HasPrefix(in, "0x") {
		in = in[2:]
	}

	return hex.DecodeString(in)
}

func MustHexToBytes(in string) []byte {
	out, err := HexToBytes(in)
	if err != nil {
		panic(err)
	}
	return out
}
