package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/tendermint/dex-demo/embedded/price/client/cli"
)

type ModuleClient struct {
	cdc *amino.Codec
}

func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{
		cdc: cdc,
	}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	priceQueryCmd := &cobra.Command{
		Use:   "price",
		Short: "queries price data",
	}
	priceQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdHistory(mc.cdc),
	)...)
	return priceQueryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	priceTxCmd := &cobra.Command{
		Use:   "price",
		Short: "manages price data",
	}
	return priceTxCmd
}
