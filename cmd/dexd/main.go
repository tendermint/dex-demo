package main

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"

	"github.com/tendermint/dex-demo/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var DefaultNodeHome = os.ExpandEnv("$HOME/.dexd")

const (
	flagOverwrite = "overwrite"
)

var mktDataDB dbm.DB

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "dexd",
		Short: "DeX daemon",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := server.PersistentPreRunEFn(ctx)(cmd, args)
			if err != nil {
				return err
			}

			mdb, err := initMktDataDB()
			if err != nil {
				return err
			}
			mktDataDB = mdb
			return nil
		},
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		client.NewCompletionCmd(rootCmd, true),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewDexApp(logger, db, mktDataDB, traceStore)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		uexApp := app.NewDexApp(logger, db, mktDataDB, traceStore)
		err := uexApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return uexApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	uexApp := app.NewDexApp(logger, db, mktDataDB, traceStore)
	return uexApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func initMktDataDB() (dbm.DB, error) {
	dir := path.Join(viper.GetString(cli.HomeFlag), "data")
	return dbm.NewGoLevelDB("mktdata", dir)
}
