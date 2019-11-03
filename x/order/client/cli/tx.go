package cli

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tendermint/dex-demo/pkg/cliutil"
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/order/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCmdStop(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [market-id] [direction] [price] [quantity] [time-in-force-blocks] [init-price] [relayer-address-hex] [relayer-fee]",
		Short: "posts an order under certain price conditions, i.e. a stop-loss or stop-buy",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}

			marketID := store.NewEntityIDFromString(args[0])
			var direction matcheng.Direction
			dirArg := strings.ToLower(args[1])
			if dirArg == "bid" {
				direction = matcheng.Bid
			} else if dirArg == "ask" {
				direction = matcheng.Ask
			} else {
				return errors.New("invalid direction")
			}

			price, err := sdk.ParseUint(args[2])
			if err != nil {
				return err
			}
			quantity, err := sdk.ParseUint(args[3])
			if err != nil {
				return err
			}
			tif, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			if tif > math.MaxUint16 {
				return errors.New("time in force too large")
			}

			initPrice, err := sdk.ParseUint(args[5])
			if err != nil {
				return err
			}
			relayedAddress, err := sdk.AccAddressFromHex(args[6])
			if err != nil {
				return err
			}
			err = sdk.VerifyAddressFormat(relayedAddress)
			if err != nil {
				return errors.New("invalid address format")
			}
			relayerFee, err := sdk.ParseCoins(args[7])
			if err != nil {
				return err
			}

			msg := types.NewMsgStop(ctx.GetFromAddress(), marketID, direction, price, quantity, uint16(tif), initPrice, relayedAddress, relayerFee)
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}

func GetCmdPost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "post [market-id] [direction] [price] [quantity] [time-in-force-blocks]",
		Short: "posts an order",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}

			marketID := store.NewEntityIDFromString(args[0])
			var direction matcheng.Direction
			dirArg := strings.ToLower(args[1])
			if dirArg == "bid" {
				direction = matcheng.Bid
			} else if dirArg == "ask" {
				direction = matcheng.Ask
			} else {
				return errors.New("invalid direction")
			}

			price, err := sdk.ParseUint(args[2])
			if err != nil {
				return err
			}
			quantity, err := sdk.ParseUint(args[3])
			if err != nil {
				return err
			}
			tif, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			if tif > math.MaxUint16 {
				return errors.New("time in force too large")
			}

			msg := types.NewMsgPost(ctx.GetFromAddress(), marketID, direction, price, quantity, uint16(tif))
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}

func GetCmdCancel(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel [order-id]",
		Short: "cancels an order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}

			orderID := store.NewEntityIDFromString(args[0])
			msg := types.NewMsgCancel(ctx.GetFromAddress(), orderID)
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}
