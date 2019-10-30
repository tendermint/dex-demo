package cliutil

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func BuildEnsuredCtx(cdc *codec.Codec) (context.CLIContext, authtypes.TxBuilder, error) {
	cliCtx := client.NewCLIContext().WithCodec(cdc)
	accGetter := authtypes.NewAccountRetriever(cliCtx)
	if err := accGetter.EnsureExists(cliCtx.GetFromAddress()); err != nil {
		return context.CLIContext{}, authtypes.TxBuilder{}, err
	}

	bldr := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	return cliCtx, bldr, nil
}

func ValidateAndBroadcast(cliCtx context.CLIContext, bldr authtypes.TxBuilder, msg sdk.Msg) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, bldr, []sdk.Msg{msg})
}
