package asset

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/dex-demo/types/errs"
	"github.com/tendermint/dex-demo/x/asset/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgMint:
			return handleMsgMint(ctx, keeper, msg)
		case types.MsgBurn:
			return handleMsgBurn(ctx, keeper, msg)
		case types.MsgTransfer:
			return handleMsgTransfer(ctx, keeper, msg)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unknown message type %v", msg.Type())).Result()
		}
	}
}

func handleMsgMint(ctx sdk.Context, keeper Keeper, msg types.MsgMint) sdk.Result {
	asset, err := keeper.Get(ctx, msg.ID)
	if err != nil {
		return err.Result()
	}
	if !asset.Owner.Equals(msg.Minter) {
		return sdk.ErrUnauthorized("cannot mint unowned asset").Result()
	}
	return errs.ErrOrBlankResult(keeper.Mint(ctx, msg.ID, msg.Amount))
}

func handleMsgBurn(ctx sdk.Context, keeper Keeper, msg types.MsgBurn) sdk.Result {
	return errs.ErrOrBlankResult(keeper.Burn(ctx, msg.ID, msg.Burner, msg.Amount))
}

func handleMsgTransfer(ctx sdk.Context, keeper Keeper, msg types.MsgTransfer) sdk.Result {
	return errs.ErrOrBlankResult(keeper.Transfer(ctx, msg.ID, msg.From, msg.To, msg.Amount))
}
