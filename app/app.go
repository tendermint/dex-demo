package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/embedded/balance"
	"github.com/tendermint/dex-demo/embedded/batch"
	"github.com/tendermint/dex-demo/embedded/book"
	"github.com/tendermint/dex-demo/embedded/fill"
	embeddedorder "github.com/tendermint/dex-demo/embedded/order"
	"github.com/tendermint/dex-demo/embedded/price"
	"github.com/tendermint/dex-demo/execution"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/x/asset"
	assettypes "github.com/tendermint/dex-demo/x/asset/types"
	"github.com/tendermint/dex-demo/x/market"
	markettypes "github.com/tendermint/dex-demo/x/market/types"
	"github.com/tendermint/dex-demo/x/order"
	ordertypes "github.com/tendermint/dex-demo/x/order/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const (
	AppName = "dex-demo"
)

var (
	// default home directories for gaiacli
	DefaultCLIHome = os.ExpandEnv("$HOME/.dexcli")

	// default home directories for gaiad
	DefaultNodeHome = os.ExpandEnv("$HOME/.dexd")

	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		asset.AppModuleBasic{},
		market.AppModuleBasic{},
		order.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc.Seal()
}

type dexApp struct {
	*baseapp.BaseApp
	Cdc *codec.Codec
	Mq  types.Backend

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// consensus keepers

	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	MintKeeper     mint.Keeper
	DistrKeeper    distr.Keeper
	ParamsKeeper   params.Keeper

	AssetKeeper  asset.Keeper
	MarketKeeper market.Keeper
	OrderKeeper  order.Keeper
	ExecKeeper   execution.Keeper

	mm *module.Manager
}

func NewDexApp(
	lgr log.Logger, appDB dbm.DB, mktDataDB dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp),
) *dexApp {
	cdc := MakeCodec()

	fillKeeper := fill.NewKeeper(mktDataDB, cdc)
	priceKeeper := price.NewKeeper(mktDataDB, cdc)
	embOrderKeeper := embeddedorder.NewKeeper(mktDataDB, cdc)
	batchKeeper := batch.NewKeeper(mktDataDB, cdc)

	queue := types.NewMemBackend()
	queue.Start()
	consumer := types.NewLocalConsumer(queue, []types.EventHandler{
		fillKeeper,
		priceKeeper,
		embOrderKeeper,
		batchKeeper,
	})
	consumer.Start()

	bApp := baseapp.NewBaseApp(AppName, lgr, appDB, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		params.StoreKey, assettypes.StoreKey, markettypes.StoreKey,
		ordertypes.StoreKey, ordertypes.LastPriceKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &dexApp{
		BaseApp: bApp,
		Mq:      queue,
		Cdc:     cdc,
		keys:    keys,
		tkeys:   tkeys,
	}

	app.ParamsKeeper = params.NewKeeper(app.Cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.ParamsKeeper.Subspace(slashing.DefaultParamspace)

	app.AccountKeeper = auth.NewAccountKeeper(app.Cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.BankKeeper = bank.NewBaseKeeper(app.AccountKeeper, bankSubspace, bank.DefaultCodespace, app.ModuleAccountAddrs())
	app.SupplyKeeper = supply.NewKeeper(app.Cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.Cdc, keys[staking.StoreKey], tkeys[staking.TStoreKey], app.SupplyKeeper, stakingSubspace, staking.DefaultCodespace,
	)
	app.MintKeeper = mint.NewKeeper(app.Cdc, keys[mint.StoreKey], mintSubspace, &stakingKeeper, app.SupplyKeeper, auth.FeeCollectorName)
	app.DistrKeeper = distr.NewKeeper(app.Cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.SupplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.SlashingKeeper = slashing.NewKeeper(
		app.Cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AssetKeeper = asset.NewKeeper(
		app.BankKeeper,
		keys[assettypes.StoreKey],
		app.Cdc,
	)
	app.MarketKeeper = market.NewKeeper(
		keys[markettypes.StoreKey],
		app.AssetKeeper,
		app.Cdc,
	)
	app.OrderKeeper = order.NewKeeper(
		app.BankKeeper,
		app.MarketKeeper,
		app.AssetKeeper,
		keys[ordertypes.StoreKey],
		keys[ordertypes.LastPriceKey],
		queue,
		app.Cdc,
	)
	app.ExecKeeper = execution.NewKeeper(
		queue,
		app.MarketKeeper,
		app.OrderKeeper,
		app.BankKeeper,
	)

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.AccountKeeper),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistrKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.DistrKeeper, app.AccountKeeper, app.SupplyKeeper),
		asset.NewAppModule(app.AssetKeeper, app.BankKeeper),
		market.NewAppModule(app.MarketKeeper),
		order.NewAppModule(app.OrderKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(staking.ModuleName)

	app.mm.SetOrderInitGenesis(
		genaccounts.ModuleName, distr.ModuleName, staking.ModuleName, auth.ModuleName,
		bank.ModuleName, slashing.ModuleName, mint.ModuleName,
		supply.ModuleName, genutil.ModuleName, assettypes.ModuleName, markettypes.ModuleName,
	)

	app.QueryRouter().
		AddRoute("embeddedorder", embeddedorder.NewQuerier(embOrderKeeper)).
		AddRoute("balance", balance.NewQuerier(app.AssetKeeper)).
		AddRoute("fill", fill.NewQuerier(fillKeeper)).
		AddRoute("price", price.NewQuerier(priceKeeper)).
		AddRoute("book", book.NewQuerier(embOrderKeeper)).
		AddRoute("batch", batch.NewQuerier(batchKeeper))

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

func (app *dexApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.Cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *dexApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *dexApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	app.performMatching(ctx)
	return app.mm.EndBlock(ctx, req)
}

func (app *dexApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *dexApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.Cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.StakingKeeper)

	return appState, validators, nil
}

func (app *dexApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])

}

func (app *dexApp) performMatching(ctx sdk.Context) {
	err := app.ExecKeeper.ExecuteAndCancelExpired(ctx)
	// an error in the execution/cancellation step is a
	// critical consensus failure.
	if err != nil {
		panic(err)
	}
}

func (app *dexApp) Codec() *codec.Codec {
	return app.Cdc
}
