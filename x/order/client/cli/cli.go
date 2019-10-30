package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(sk string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "order",
		Short: "queries orders",
	}
	queryCmd.AddCommand(client.GetCommands(
		GetCmdListOrders(sk, cdc),
	)...)
	return queryCmd
}

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "order",
		Short: "manages orders",
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdPost(cdc),
		GetCmdCancel(cdc),
	)...)
	return txCmd
}
