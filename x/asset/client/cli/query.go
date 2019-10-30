package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/dex-demo/x/asset/types"
)

func GetCmdListAssets(queryRoute string, cdc *codec.Codec) *cobra.Command {
	var symbol string
	out := &cobra.Command{
		Use:   "list",
		Short: "lists assets with the specified filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/list/%s", queryRoute, symbol), nil)
			if err != nil {
				return err
			}

			var out types.ListQueryResult
			cdc.MustUnmarshalJSON(res, &out)
			return ctx.PrintOutput(out)
		},
	}
	out.Flags().StringVar(&symbol, "symbol", "", "filter by symbol")
	return out
}
