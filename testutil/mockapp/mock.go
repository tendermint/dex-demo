package mockapp

import (
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/app"
	"github.com/tendermint/dex-demo/execution"
	"github.com/tendermint/dex-demo/types"
	"github.com/tendermint/dex-demo/x/asset"
	"github.com/tendermint/dex-demo/x/market"
	"github.com/tendermint/dex-demo/x/order"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type MockApp struct {
	Cdc             *codec.Codec
	Mq              types.Backend
	Ctx             sdk.Context
	AssetKeeper     asset.Keeper
	MarketKeeper    market.Keeper
	OrderKeeper     order.Keeper
	BankKeeper      bank.Keeper
	ExecutionKeeper execution.Keeper
}

type Option func(t *testing.T, app *MockApp)

func New(t *testing.T, options ...Option) *MockApp {
	appDB := dbm.NewMemDB()
	mkDataDB := dbm.NewMemDB()
	dex := app.NewDexApp(log.NewNopLogger(), appDB, mkDataDB, &nopWriter{})
	dex.InitChain(abci.RequestInitChain{
		AppStateBytes: []byte("{}"),
	})
	ctx := dex.BaseApp.NewContext(false, abci.Header{ChainID: "unit-test-chain", Height: 1, Time: time.Unix(1558332092, 0)})

	mock := &MockApp{
		Cdc:             dex.Cdc,
		Mq:              dex.Mq,
		Ctx:             ctx,
		AssetKeeper:     dex.AssetKeeper,
		MarketKeeper:    dex.MarketKeeper,
		OrderKeeper:     dex.OrderKeeper,
		BankKeeper:      dex.BankKeeper,
		ExecutionKeeper: dex.ExecKeeper,
	}

	for _, opt := range options {
		opt(t, mock)
	}

	return mock
}
