package order

import (
	"net/http"

	"github.com/tendermint/dex-demo/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/embedded/auth"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/user/orders", auth.DefaultAuthMW(getOrdersHandler(ctx, cdc))).Methods("GET")
}

func getOrdersHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := auth.MustGetKBFromSession(r)
		q := r.URL.Query()

		req := ListQueryRequest{
			Owner: owner.GetAddr(),
		}
		if start, ok := q["start"]; ok {
			req.Start = sdk.NewUintFromString(start[0])
		}

		resB, _, err := ctx.QueryWithData("custom/embeddedorder/list", cdc.MustMarshalBinaryBare(req))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		embedded.PostProcessResponse(w, ctx, resB)
	}
}
