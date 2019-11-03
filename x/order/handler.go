package order

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/pkg/log"
	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/x/order/types"
)

var logger = log.WithModule("order")

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgStop:
			return handleMsgStop(ctx, keeper, msg)
		case types.MsgPost:
			return handleMsgPost(ctx, keeper, msg)
		case types.MsgCancel:
			return handleMsgCancel(ctx, keeper, msg)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unknown message type %v", msg.Type())).Result()
		}
	}
}

func handleMsgStop(ctx sdk.Context, keeper Keeper, msg types.MsgStop) sdk.Result {
	currentPrice := sdk.NewUint(0)
	// Shouldn't be triggered in these ranges
	if msg.Post.Price.GT(msg.PastPrice) && currentPrice.LT(msg.Post.Price) {
		return sdk.Result{
			Log: fmt.Sprintf("current price not in triggering range"),
		}
	}
	if msg.Post.Price.LT(msg.PastPrice) && currentPrice.GT(msg.Post.Price) {
		return sdk.Result{
			Log: fmt.Sprintf("current price not in triggering range"),
		}
	}

	order, err := keeper.Post(
		ctx,
		msg.Post.Owner,
		msg.Post.MarketID,
		msg.Post.Direction,
		msg.Post.Price,
		msg.Post.Quantity,
		msg.Post.TimeInForce,
	)
	if err != nil {
		return err.Result()
	}
	logger.Info(
		"stop order",
		"id", order.ID.String(),
		"market_id", order.MarketID.String(),
		"price", order.Price.String(),
		"quantity", order.Quantity.String(),
		"direction", order.Direction.String(),
	)

	err = keeper.bankKeeper.SendCoins(ctx, msg.Post.Owner, msg.Relayer, msg.RelayFee)
	if err == nil {
		return err.Result()
	}
	logger.Info(
		"stop order",
		"id", order.ID.String(),
		"market_id", order.MarketID.String(),
		"price", order.Price.String(),
		"quantity", order.Quantity.String(),
		"direction", order.Direction.String(),
	)

	return sdk.Result{
		Log: fmt.Sprintf("order_id:%s", order.ID),
	}
}

func handleMsgPost(ctx sdk.Context, keeper Keeper, msg types.MsgPost) sdk.Result {
	order, err := keeper.Post(
		ctx,
		msg.Owner,
		msg.MarketID,
		msg.Direction,
		msg.Price,
		msg.Quantity,
		msg.TimeInForce,
	)

	if err == nil {
		logger.Info(
			"posted order",
			"id", order.ID.String(),
			"market_id", order.MarketID.String(),
			"price", order.Price.String(),
			"quantity", order.Quantity.String(),
			"direction", order.Direction.String(),
		)
		return sdk.Result{
			Log: fmt.Sprintf("order_id:%s", order.ID),
		}
	}

	return err.Result()
}

func handleMsgCancel(ctx sdk.Context, keeper Keeper, msg types.MsgCancel) sdk.Result {
	order, err := keeper.Get(ctx, msg.OrderID)
	if err != nil {
		return err.Result()
	}
	if !order.Owner.Equals(msg.Owner) {
		return sdk.ErrUnauthorized("cannot cancel unowned order").Result()
	}
	return errs.ErrOrBlankResult(keeper.Del(ctx, order.ID))
}
