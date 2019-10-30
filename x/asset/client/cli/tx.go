package cli

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/dex-demo/pkg/cliutil"
	"github.com/tendermint/dex-demo/types/store"
	"github.com/tendermint/dex-demo/x/asset/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCmdMint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint [asset id] [amount]",
		Short: "mints additional assets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}

			msg := types.NewMsgMint(store.NewEntityIDFromString(args[0]), ctx.GetFromAddress(), sdk.NewUintFromString(args[1]))
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}

func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [asset id] [amount]",
		Short: "burns assets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}
			msg := types.NewMsgBurn(store.NewEntityIDFromString(args[0]), ctx.GetFromAddress(), sdk.NewUintFromString(args[1]))
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}

func GetCmdTransfer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [to address] [asset id] [amount]",
		Short: "transfers assets",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, bldr, err := cliutil.BuildEnsuredCtx(cdc)
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransfer(store.NewEntityIDFromString(args[1]), ctx.GetFromAddress(), toAddr, sdk.NewUintFromString(args[2]))
			return cliutil.ValidateAndBroadcast(ctx, bldr, msg)
		},
	}
}
