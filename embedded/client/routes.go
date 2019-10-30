package client

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/dex-demo/embedded/auth"
	"github.com/tendermint/dex-demo/embedded/balance"
	"github.com/tendermint/dex-demo/embedded/batch"
	"github.com/tendermint/dex-demo/embedded/book"
	"github.com/tendermint/dex-demo/embedded/exchange"
	"github.com/tendermint/dex-demo/embedded/fill"
	"github.com/tendermint/dex-demo/embedded/order"
	"github.com/tendermint/dex-demo/embedded/price"
	"github.com/tendermint/dex-demo/embedded/ui"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec, enableFaucet bool) {
	r.Use(auth.HandleCORSMW)
	r.Use(auth.ProtectCSRFMW([]string{
		"/api/v1/faucet/transfer",
	}))
	sub := r.PathPrefix("/api/v1").Subrouter()
	auth.RegisterRoutes(ctx, sub, cdc)
	exchange.RegisterRoutes(ctx, sub, cdc)
	fill.RegisterRoutes(ctx, sub, cdc)
	order.RegisterRoutes(ctx, sub, cdc)
	balance.RegisterRoutes(ctx, sub, cdc, enableFaucet)
	price.RegisterRoutes(ctx, sub, cdc)
	book.RegisterRoutes(ctx, sub, cdc)
	batch.RegisterRoutes(ctx, sub, cdc)
	ui.RegisterRoutes(ctx, r, cdc)
}
