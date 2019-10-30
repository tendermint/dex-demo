package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(sk string, cdc *codec.Codec) *cobra.Command {
	marketQueryCmd := &cobra.Command{
		Use:   "market",
		Short: "queries available markets",
	}
	marketQueryCmd.AddCommand(client.GetCommands(
		GetCmdListMarkets(sk, cdc),
	)...)
	return marketQueryCmd
}

func GetTxCmd() *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:   "market",
		Short: "manages available markets",
	}
	return marketTxCmd
}
