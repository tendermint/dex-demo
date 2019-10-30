package matcheng

import (
	"encoding/json"
	"errors"
)

const (
	Bid Direction = iota
	Ask
)

type Direction uint8

func (d Direction) String() string {
	if d == Bid {
		return "BID"
	}

	return "ASK"
}

func (d *Direction) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	if str == "BID" {
		*d = Bid
	} else if str == "ASK" {
		*d = Ask
	} else {
		return errors.New("invalid direction")
	}

	return nil
}

func (d Direction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}
