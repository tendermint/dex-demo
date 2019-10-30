package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(sk string, cdc *codec.Codec) *cobra.Command {
	assetQueryCmd := &cobra.Command{
		Use:   "asset",
		Short: "queries on-chain assets",
	}
	assetQueryCmd.AddCommand(client.GetCommands(
		GetCmdListAssets(sk, cdc),
	)...)
	return assetQueryCmd
}

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "asset",
		Short: "manages on-chain assets",
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdMint(cdc),
		GetCmdBurn(cdc),
		GetCmdTransfer(cdc),
	)...)
	return txCmd
}
