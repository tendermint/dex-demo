package execution

import (
	"time"

	"github.com/tendermint/dex-demo/pkg/log"
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/types/store"
	assettypes "github.com/tendermint/dex-demo/x/asset/types"
	"github.com/tendermint/dex-demo/x/market"
	"github.com/tendermint/dex-demo/x/order"
	types2 "github.com/tendermint/dex-demo/x/order/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	queue     types.Backend
	mk        market.Keeper
	ordK      order.Keeper
	bk        bank.Keeper
	metrics   *Metrics
	saveFills bool
}

type matcherByMarket struct {
	matcher *matcheng.Matcher
	mktID   store.EntityID
}

var logger = log.WithModule("execution")

func NewKeeper(queue types.Backend, mk market.Keeper, ordK order.Keeper, bk bank.Keeper) Keeper {
	return Keeper{
		queue:   queue,
		mk:      mk,
		ordK:    ordK,
		bk:      bk,
		metrics: PrometheusMetrics(),
	}
}

func (k Keeper) ExecuteAndCancelExpired(ctx sdk.Context) sdk.Error {
	start := time.Now()
	height := ctx.BlockHeight()

	var toCancel []store.EntityID
	k.ordK.Iterator(ctx, func(ord types2.Order) bool {
		if height-ord.CreatedBlock > int64(ord.TimeInForceBlocks) {
			toCancel = append(toCancel, ord.ID)
		}

		return true
	})
	for _, ordID := range toCancel {
		if err := k.ordK.Cancel(ctx, ordID); err != nil {
			return err
		}
	}

	logger.Info("cancelled expired orders", "count", len(toCancel))

	matchersByMarket := make(map[string]*matcherByMarket)

	k.ordK.ReverseIterator(ctx, func(ord types2.Order) bool {
		matcher := getMatcherByMarket(matchersByMarket, ord).matcher
		matcher.EnqueueOrder(ord.Direction, ord.ID, ord.Price, ord.Quantity)
		return true
	})

	var toFill []*matcheng.MatchResults
	for _, m := range matchersByMarket {
		res := m.matcher.Match()
		if res == nil {
			continue
		}

		_ = k.queue.Publish(types.Batch{
			BlockNumber:   height,
			BlockTime:     ctx.BlockHeader().Time,
			MarketID:      m.mktID,
			ClearingPrice: res.ClearingPrice,
			Bids:          res.BidAggregates,
			Asks:          res.AskAggregates,
		})
		toFill = append(toFill, res)
		matcheng.ReturnMatcher(m.matcher)
	}
	var fillCount int
	for _, res := range toFill {
		fillCount += len(res.Fills)
		for _, f := range res.Fills {
			if err := k.ExecuteFill(ctx, res.ClearingPrice, f); err != nil {
				return err
			}
		}
	}

	logger.Info("matched orders", "count", fillCount)

	duration := time.Since(start).Nanoseconds()
	k.metrics.ProcessingTime.Observe(float64(duration) / 1000000)
	k.metrics.OrdersProcessed.Observe(float64(fillCount))
	return nil
}

func (k Keeper) ExecuteFill(ctx sdk.Context, clearingPrice sdk.Uint, f matcheng.Fill) sdk.Error {
	ord, err := k.ordK.Get(ctx, f.OrderID)
	if err != nil {
		return err
	}
	mkt, err := k.mk.Get(ctx, ord.MarketID)
	if err != nil {
		return err
	}
	pair, err := k.mk.Pair(ctx, mkt.ID)
	if err != nil {
		panic(err)
	}

	if ord.Direction == matcheng.Bid {
		quoteAmount := f.QtyFilled
		_, err = k.bk.AddCoins(ctx, ord.Owner, assettypes.Coins(mkt.BaseAssetID, quoteAmount))
		if err != nil {
			return err
		}
		if clearingPrice.LT(ord.Price) {
			diff := ord.Price.Sub(clearingPrice)
			refund, qErr := matcheng.NormalizeQuoteQuantity(diff, f.QtyFilled)
			if qErr == nil {
				_, err = k.bk.AddCoins(ctx, ord.Owner, assettypes.Coins(mkt.QuoteAssetID, refund))
				if err != nil {
					return err
				}
			} else {
				logger.Info(
					"refund amount too small",
					"order_id", ord.ID.String(),
					"qty_filled", f.QtyFilled.String(),
					"price_delta", diff.String(),
				)
			}
		}
	} else {
		baseAmount, qErr := matcheng.NormalizeQuoteQuantity(clearingPrice, f.QtyFilled)
		if qErr == nil {
			_, err = k.bk.AddCoins(ctx, ord.Owner, assettypes.Coins(mkt.QuoteAssetID, baseAmount))
			if err != nil {
				return err
			}
		} else {
			panic("clearing price too small to represent")
		}
	}

	ord.Quantity = f.QtyUnfilled
	if ord.Quantity.Equal(sdk.ZeroUint()) {
		logger.Info("order completely filled", "id", ord.ID.String())
		if err := k.ordK.Del(ctx, ord.ID); err != nil {
			return err
		}
	} else {
		logger.Info("order partially filled", "id", ord.ID.String())
		if err := k.ordK.Set(ctx, ord); err != nil {
			return err
		}
	}

	_ = k.queue.Publish(types.Fill{
		OrderID:     ord.ID,
		MarketID:    mkt.ID,
		Owner:       ord.Owner,
		Pair:        pair,
		Direction:   ord.Direction,
		QtyFilled:   f.QtyFilled,
		QtyUnfilled: f.QtyUnfilled,
		BlockNumber: ctx.BlockHeight(),
		BlockTime:   ctx.BlockHeader().Time.Unix(),
		Price:       clearingPrice,
	})
	return nil
}

func getMatcherByMarket(matchers map[string]*matcherByMarket, ord types2.Order) *matcherByMarket {
	mkt := ord.MarketID.String()
	matcher := matchers[mkt]
	if matcher == nil {
		matcher = &matcherByMarket{
			matcher: matcheng.GetMatcher(),
			mktID:   ord.MarketID,
		}
		matchers[mkt] = matcher
	}
	return matcher
}
