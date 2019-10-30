package serde

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type HexBytes []byte

func (h *HexBytes) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		*h = nil
		return nil
	}
	if len(buf) <= 2 {
		return errors.New("no value")
	}

	unquoted := string(buf[1 : len(buf)-1])
	data, err := hex.DecodeString(unquoted[2:])
	if err != nil {
		return err
	}
	*h = data
	return nil
}

func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(fmt.Sprintf("0x%s", hex.EncodeToString(h)))
}

func MustMarshalSortedJSON(in interface{}) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
