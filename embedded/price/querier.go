package price

import (
	"time"

	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/types/errs"
)

const (
	QueryHistory = "history"
	QueryCandles = "candles"
	QueryDaily   = "daily"
	MaxTicks     = 2000
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryHistory:
			return queryHistory(path[1:], keeper)
		case QueryCandles:
			return queryCandles(path[1:], req.Data, keeper)
		case QueryDaily:
			return queryDaily(path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown price query endpoint")
		}
	}
}

func queryHistory(path []string, keeper Keeper) ([]byte, sdk.Error) {
	mktID := sdk.NewUintFromString(path[0])

	res := TickQueryResult{
		MarketID: mktID,
		Ticks:    make([]TickEntry, 0),
	}

	keeper.ReverseIteratorByMarket(mktID, func(tick Tick) bool {
		if res.Pair == "" {
			res.Pair = tick.Pair
		}

		res.Ticks = append(res.Ticks, TickEntry{
			BlockNumber: tick.BlockNumber,
			Timestamp:   tick.BlockTime,
			Price:       tick.Price,
		})

		return len(res.Ticks) < MaxTicks
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, errs.ErrMarshalFailure("could not marshal tick result")
	}
	return b, nil
}

func queryCandles(path []string, data []byte, keeper Keeper) ([]byte, sdk.Error) {
	mktID := sdk.NewUintFromString(path[0])

	var params CandleQueryParams
	err := amino.UnmarshalBinaryBare(data, &params)
	if err != nil {
		return nil, errs.ErrUnmarshalFailure("could not unmarshal query params")
	}
	if params.From.After(params.To) {
		return nil, errs.ErrInvalidArgument("from cannot be after to")
	}

	res := CandleQueryResult{
		MarketID: mktID,
	}

	delta := params.Interval.Delta()
	keeper.IteratorByMarketAndInterval(mktID, params.From, params.To, func(tick Tick) bool {
		if res.Pair == "" {
			res.Pair = tick.Pair
		}
		if len(res.Candles) == 0 {
			res.Candles = append(res.Candles, CandleEntry{
				Date:  roundTime(time.Unix(tick.BlockTime, 0), params.Interval),
				Open:  tick.Price,
				High:  tick.Price,
				Low:   tick.Price,
				Close: tick.Price,
			})
		}
		lastCandle := &res.Candles[len(res.Candles)-1]
		if tick.BlockTime-lastCandle.Date.Unix() > delta {
			res.Candles = append(res.Candles, CandleEntry{
				Date:  roundTime(time.Unix(tick.BlockTime, 0), params.Interval),
				Open:  tick.Price,
				High:  tick.Price,
				Low:   tick.Price,
				Close: tick.Price,
			})
			lastCandle = &res.Candles[len(res.Candles)-1]
		}
		lastCandle.Close = tick.Price
		if tick.Price.GT(lastCandle.High) {
			lastCandle.High = tick.Price
		}
		if tick.Price.LT(lastCandle.Low) {
			lastCandle.Low = tick.Price
		}

		return len(res.Candles) < MaxTicks
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, errs.ErrMarshalFailure("could not marshal candle result")
	}
	return b, nil
}

func queryDaily(path []string, keeper Keeper) ([]byte, sdk.Error) {
	mktID := sdk.NewUintFromString(path[0])

	res := DailyQueryResult{
		Pair:   "",
		Volume: sdk.ZeroUint(),
		Change: sdk.ZeroDec(),
		Last:   sdk.ZeroUint(),
		High:   sdk.ZeroUint(),
		Low:    sdk.ZeroUint(),
	}

	now := time.Now()
	startTime := now.Add(time.Duration(-24) * time.Hour)
	keeper.IteratorByMarketAndInterval(mktID, startTime, now, func(tick Tick) bool {
		if res.Pair == "" {
			res.Pair = tick.Pair
		}
		if res.Last.IsZero() {
			res.Last = tick.Price
		}
		if res.High.LT(tick.Price) {
			res.High = tick.Price
		}
		if res.Low.IsZero() || res.Low.GT(tick.Price) {
			res.Low = tick.Price
		}
		return true
	})
	if res.Pair == "" {
		return nil, sdk.ErrInternal("no price points found")
	}

	prevClose := sdk.ZeroUint()
	keeper.ReverseIteratorByMarketFrom(mktID, startTime, func(tick Tick) bool {
		if tick.BlockTime == startTime.Unix() {
			return true
		}

		prevClose = tick.Price
		return false
	})

	if prevClose.IsZero() {
		res.Change = sdk.OneDec()
	} else {
		res.Change = sdk.NewDecFromBigInt(res.Last.BigInt()).Quo(sdk.NewDecFromBigInt(prevClose.BigInt()))
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result")
	}
	return b, nil
}

func roundTime(t time.Time, interval CandleInterval) time.Time {
	switch interval {
	case CandleInterval1M:
		return t.Truncate(time.Minute)
	case CandleInterval5M:
		return t.Truncate(5 * time.Minute)
	case CandleInterval15M:
		return t.Truncate(15 * time.Minute)
	case CandleInterval30M:
		return t.Truncate(30 * time.Minute)
	case CandleInterval60M:
		return t.Truncate(60 * time.Minute)
	default:
		panic("invalid time interval")
	}
}
