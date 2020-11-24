package matcheng

import (
	"sort"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/log"
)

type Order struct {
	ID       sdk.Uint
	Price    sdk.Uint
	Quantity sdk.Uint
}

type MatchResults struct {
	ClearingPrice sdk.Uint
	Fills         []Fill
	MatchTable    []MatchEntry
	BidAggregates []AggregatePrice
	AskAggregates []AggregatePrice
}

type Matcher struct {
	bids []Order
	asks []Order

	mtx sync.Mutex
}

var logger = log.WithModule("matcher")

type MatchEntry [3]sdk.Uint
type AggregatePrice [2]sdk.Uint

var zero = sdk.ZeroUint()

func NewMatcher() *Matcher {
	return &Matcher{
		bids: make([]Order, 0),
		asks: make([]Order, 0),
	}
}

// merge prices into one big list
// iterate over prices
// choose the price at with # of sells > # of buys
// walk back one price - that's clearing.
// degen case: vertical line (then choose midpoint)
// other degen case: no overlap.

func (m *Matcher) EnqueueOrder(oType Direction, id sdk.Uint, price sdk.Uint, quantity sdk.Uint) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	order := &Order{
		ID:       id,
		Price:    price,
		Quantity: quantity,
	}

	if oType == Bid {
		m.enqueueBid(order)
	} else {
		m.enqueueAsk(order)
	}
}

func (m *Matcher) Match() *MatchResults {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.bids) == 0 || len(m.asks) == 0 {
		logger.Info("no bids or asks in this block")
		return nil
	}

	// handle degenerate case 1: no matches
	if m.bids[len(m.bids)-1].Price.LT(m.asks[0].Price) {
		logger.Info("highest bid price is lower than lowest ask price")
		return nil
	}

	// [price, supply, demand]
	var matchTable []MatchEntry
	var askAggs []AggregatePrice
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		i := 0
		for _, ask := range m.asks {
			if len(matchTable) == 0 {
				matchTable = append(matchTable, MatchEntry{ask.Price, ask.Quantity, zero})
				askAggs = append(askAggs, AggregatePrice{ask.Price, ask.Quantity})

				i++
				continue
			}

			last := len(matchTable) - 1
			if matchTable[last][0].Equal(ask.Price) {
				matchTable[last][1] = matchTable[last][1].Add(ask.Quantity)
				askAggs[last][1] = askAggs[last][1].Add(ask.Quantity)
				continue
			}

			matchTable = append(matchTable, MatchEntry{
				ask.Price,
				matchTable[i-1][1].Add(ask.Quantity),
				zero,
			})
			askAggs = append(askAggs, AggregatePrice{
				ask.Price,
				matchTable[i-1][1].Add(ask.Quantity),
			})
			i++
		}

		wg.Done()
	}()

	var bidAggs []AggregatePrice
	go func() {
		for i := len(m.bids) - 1; i >= 0; i-- {
			bid := m.bids[i]

			if len(bidAggs) == 0 {
				bidAggs = append(bidAggs, [2]sdk.Uint{
					bid.Price,
					bid.Quantity,
				})
				continue
			}

			if bidAggs[0][0].Equal(bid.Price) {
				bidAggs[0][1] = bidAggs[0][1].Add(bid.Quantity)
				continue
			}

			bidAggs = append([]AggregatePrice{
				{bid.Price, bidAggs[0][1].Add(bid.Quantity)},
			}, bidAggs...)
		}
		wg.Done()
	}()

	wg.Wait()

	var lastInsertion int
	for _, agg := range bidAggs {
		j := sort.Search(len(matchTable), func(i int) bool {
			return matchTable[i][0].GTE(agg[0])
		})

		if j == len(matchTable) {
			carry := zero
			if len(matchTable) > 0 {
				carry = matchTable[j-1][1]
			}

			matchTable = append(matchTable, MatchEntry{agg[0], carry, agg[1]})
		} else if matchTable[j][0].Equal(agg[0]) {
			matchTable[j][2] = matchTable[j][2].Add(agg[1])
		} else {
			curr := zero
			if j > 0 {
				curr = matchTable[j-1][1]
			}

			matchTable = append(matchTable, MatchEntry{})
			copy(matchTable[j+1:], matchTable[j:])
			matchTable[j] = MatchEntry{agg[0], curr, agg[1]}
		}

		for i := lastInsertion; i < j; i++ {
			matchTable[i][2] = matchTable[i][2].Add(agg[1])
		}
		lastInsertion = j + 1
	}

	clearingPrice := zero
	aggSupply := zero
	aggDemand := zero
	crossPoint := 0
	for i, entry := range matchTable {
		crossed := i > 0 && calcDir(matchTable[i-1]) != calcDir(matchTable[i])
		if crossed {
			crossSup, crossDem := entry[1], entry[2]
			var topIdx int
			for j := i + 1; j < len(matchTable); j++ {
				testEntry := matchTable[j]
				if !crossSup.Equal(testEntry[1]) || !crossDem.Equal(testEntry[2]) {
					break
				}
				topIdx = j
			}

			if topIdx != 0 {
				diff := matchTable[topIdx][0].Sub(matchTable[i][0])
				mid := matchTable[i][0].Add(diff.Quo(sdk.NewUint(2)))

				clearingPrice = mid
				aggSupply = crossSup
				aggDemand = crossDem
				break
			}

			crossPoint = i
		}

		if i > 0 && matchTable[i-1][1].GT(zero) && crossed {
			break
		}

		clearingPrice = entry[0]
		aggSupply = entry[1]
		aggDemand = entry[2]
	}

	// check for other edge case: horizontal cross
	if clearingPrice.Equal(matchTable[len(matchTable)-1][0]) && crossPoint < len(matchTable)-1 {
		entry := matchTable[crossPoint]
		if entry[1].GT(sdk.ZeroUint()) && entry[2].GT(sdk.ZeroUint()) {
			clearingPrice = entry[0]
			aggSupply = entry[1]
			aggDemand = entry[2]
		}
	}

	aggDemandDec := sdk.NewDecFromBigInt(aggDemand.BigInt())
	aggSupplyDec := sdk.NewDecFromBigInt(aggSupply.BigInt())

	proRataDec := aggSupplyDec.Quo(aggDemandDec)
	proRataRecip := sdk.OneDec().Quo(proRataDec)
	overOne := proRataDec.GT(sdk.OneDec())
	maxBidVolume := sdk.NewDecFromInt(aggDemandDec.Mul(proRataDec).RoundInt())
	maxAskVolume := sdk.NewDecFromInt(aggSupplyDec.Mul(proRataRecip).RoundInt())

	matchedBidVolume := sdk.ZeroDec()
	var fills []Fill
	for i := len(m.bids) - 1; i >= 0; i-- {
		bid := m.bids[i]
		if bid.Price.LT(clearingPrice) || matchedBidVolume.Equal(maxAskVolume) {
			break
		}

		var qtyInt sdk.Uint
		if overOne {
			qtyInt = bid.Quantity
		} else {
			qtyDec := sdk.NewDecFromBigInt(bid.Quantity.BigInt()).Mul(proRataDec).Ceil()
			qtyInt = sdk.NewUintFromString(qtyDec.RoundInt().String())

			if matchedBidVolume.Add(qtyDec).GT(maxBidVolume) {
				qtyDec = maxBidVolume.Sub(matchedBidVolume)
				qtyInt = sdk.NewUintFromString(qtyDec.RoundInt().String())
			}

			matchedBidVolume = matchedBidVolume.Add(qtyDec)

			if qtyInt.IsZero() {
				continue
			}
		}

		fills = append(fills, Fill{
			OrderID:     bid.ID,
			QtyFilled:   qtyInt,
			QtyUnfilled: bid.Quantity.Sub(qtyInt),
		})
	}

	matchedAskVolume := sdk.ZeroDec()
	for _, ask := range m.asks {
		if ask.Price.GT(clearingPrice) || matchedAskVolume.Equal(maxAskVolume) {
			break
		}

		var qtyInt sdk.Uint
		if overOne {
			qtyDec := proRataRecip.Mul(sdk.NewDecFromBigInt(ask.Quantity.BigInt())).Ceil()
			qtyInt = sdk.NewUintFromString(qtyDec.RoundInt().String())

			if matchedAskVolume.Add(qtyDec).GT(maxAskVolume) {
				qtyDec = maxAskVolume.Sub(matchedAskVolume)
				qtyInt = sdk.NewUintFromString(qtyDec.RoundInt().String())
			}

			matchedAskVolume = matchedAskVolume.Add(qtyDec)

			if qtyInt.IsZero() {
				continue
			}
		} else {
			qtyInt = ask.Quantity
		}

		fills = append(fills, Fill{
			OrderID:     ask.ID,
			QtyFilled:   qtyInt,
			QtyUnfilled: ask.Quantity.Sub(qtyInt),
		})
	}

	logger.Info(
		"generated match results",
		"fill_count", len(fills),
		"bid_count", len(m.bids),
		"ask_count", len(m.asks),
	)

	return &MatchResults{
		ClearingPrice: clearingPrice,
		Fills:         fills,
		MatchTable:    matchTable,
		BidAggregates: bidAggs,
		AskAggregates: askAggs,
	}
}

func (m *Matcher) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.bids = make([]Order, 0)
	m.asks = make([]Order, 0)
}

// lowest bid (buy price) first

func (m *Matcher) enqueueBid(order *Order) {
	if len(m.bids) == 0 {
		m.bids = append(m.bids, *order)
		return
	}

	i := sort.Search(len(m.bids), func(i int) bool {
		tester := m.bids[i]
		if tester.Price.Equal(order.Price) {
			return tester.ID.LT(order.ID)
		}
		return tester.Price.GT(order.Price)
	})

	m.bids = append(m.bids, Order{})
	copy(m.bids[i+1:], m.bids[i:])
	m.bids[i] = *order
}

// lowest ask (sell price) first

func (m *Matcher) enqueueAsk(order *Order) {
	if len(m.asks) == 0 {
		m.asks = append(m.asks, *order)
		return
	}

	i := sort.Search(len(m.asks), func(i int) bool {
		tester := m.asks[i]
		if tester.Price.Equal(order.Price) {
			return tester.ID.LT(order.ID)
		}
		return tester.Price.GT(order.Price)
	})

	m.asks = append(m.asks, Order{})
	copy(m.asks[i+1:], m.asks[i:])
	m.asks[i] = *order
}

func calcDir(entry [3]sdk.Uint) int {
	if entry[1].GT(entry[2]) {
		return 1
	} else if entry[1].Equal(entry[2]) {
		return 0
	}
	return -1
}
