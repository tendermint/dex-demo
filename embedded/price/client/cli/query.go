package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/dex-demo/embedded/price"
)

func GetCmdHistory(cdc *codec.Codec) *cobra.Command {
	var plot bool
	getCmd := &cobra.Command{
		Use:   "history [market id]",
		Short: "get historical prices for the provided market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/history/%s", price.EntityName, args[0]), nil)
			if err != nil {
				return err
			}

			var out price.TickQueryResult
			cdc.MustUnmarshalJSON(res, &out)

			if !plot {
				return ctx.PrintOutput(out)
			}

			bin, err := exec.LookPath("gnuplot")
			if err != nil {
				return errors.New("gnuplot not found")
			}
			tmp, err := ioutil.TempFile("", "plot")
			if err != nil {
				return err
			}

			for _, tick := range out.Ticks {
				entry := fmt.Sprintf("%s %s\n", strconv.Itoa(int(tick.Timestamp)), tick.Price)
				if _, err := tmp.Write([]byte(entry)); err != nil {
					return err
				}
			}

			plotScript := fmt.Sprintf("datafile='%s';set xlabel 'Date';set ylabel 'Price';set xdata time;set timefmt '%%s';set format x '%%m/%%d/%%Y %%H:%%M:%%S';set xtics rotate;plot datafile using 1:2 with lines title '%s'", tmp.Name(), out.Pair)
			plotCmd := exec.Command(
				bin,
				"-p",
				"-e",
				plotScript,
			)
			return plotCmd.Run()
		},
	}
	getCmd.Flags().BoolVar(&plot, "plot", false, "plot the prices using gnuplot")
	return getCmd
}
