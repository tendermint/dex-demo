package price

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Tick struct {
	MarketID    sdk.Uint
	Pair        string
	BlockNumber int64
	BlockTime   int64
	Price       sdk.Uint
}

type TickEntry struct {
	BlockNumber int64    `json:"block_number"`
	Timestamp   int64    `json:"timestamp"`
	Price       sdk.Uint `json:"price"`
}

type TickQueryResult struct {
	MarketID sdk.Uint    `json:"market_id"`
	Pair     string      `json:"pair"`
	Ticks    []TickEntry `json:"ticks"`
}

func (t TickQueryResult) String() string {
	var buf strings.Builder
	table := tablewriter.NewWriter(&buf)
	table.SetCaption(true, fmt.Sprintf("%s (ID: %s)", t.Pair, t.MarketID))
	table.SetHeader([]string{
		"Block Number",
		"Timestamp",
		"Price",
	})
	for _, entry := range t.Ticks {
		table.Append([]string{
			strconv.Itoa(int(entry.BlockNumber)),
			time.Unix(entry.Timestamp, 0).String(),
			entry.Price.String(),
		})
	}
	table.Render()
	return buf.String()
}

type CandleInterval string

const (
	CandleInterval1M  CandleInterval = "1m"
	CandleInterval5M                 = "5m"
	CandleInterval15M                = "15m"
	CandleInterval30M                = "30m"
	CandleInterval60M                = "60m"
)

var validIntervals = map[string]CandleInterval{
	"1m":  CandleInterval1M,
	"5m":  CandleInterval5M,
	"15m": CandleInterval15M,
	"30m": CandleInterval30M,
	"60m": CandleInterval60M,
}

func NewCandleIntervalFromString(in string) (CandleInterval, error) {
	val, ok := validIntervals[in]
	if !ok {
		return CandleInterval1M, errors.New("unknown candle interval")
	}

	return val, nil
}

func (c CandleInterval) Delta() int64 {
	switch c {
	case CandleInterval1M:
		return 60
	case CandleInterval5M:
		return 300
	case CandleInterval15M:
		return 900
	case CandleInterval30M:
		return 1800
	case CandleInterval60M:
		return 3600
	default:
		panic("invalid candle interval")
	}
}

func (c *CandleInterval) UnmarshalJSON(data []byte) error {
	var valStr string
	err := json.Unmarshal(data, &valStr)
	if err != nil {
		return err
	}

	val, ok := validIntervals[valStr]
	if !ok {
		return errors.New("unknown candle interval")
	}

	*c = val
	return nil
}

type CandleQueryParams struct {
	From     time.Time
	To       time.Time
	Interval CandleInterval
}

type CandleQueryResult struct {
	MarketID sdk.Uint      `json:"market_id"`
	Pair     string        `json:"pair"`
	Candles  []CandleEntry `json:"candles"`
}

type CandleEntry struct {
	Date  time.Time `json:"date"`
	Open  sdk.Uint  `json:"open"`
	Close sdk.Uint  `json:"close"`
	High  sdk.Uint  `json:"high"`
	Low   sdk.Uint  `json:"low"`
}

type DailyQueryResult struct {
	Pair   string   `json:"pair"`
	Volume sdk.Uint `json:"volume"`
	Change sdk.Dec  `json:"change"`
	Last   sdk.Uint `json:"last"`
	High   sdk.Uint `json:"high"`
	Low    sdk.Uint `json:"low"`
}
