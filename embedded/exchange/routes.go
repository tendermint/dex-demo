package exchange

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/tendermint/dex-demo/embedded"
	"github.com/tendermint/dex-demo/embedded/auth"
	"github.com/tendermint/dex-demo/x/order/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	sdkauth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	sub := r.PathPrefix("/exchange").Subrouter()
	sub.Use(auth.DefaultAuthMW)
	sub.HandleFunc("/orders", postOrderHandler(ctx, cdc)).Methods("POST")
}

func postOrderHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OrderCreationRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		kb := auth.MustGetKBFromSession(r)
		owner := kb.GetAddr()
		ctx = ctx.WithFromAddress(owner)

		msg := types.NewMsgPost(owner, req.MarketID, req.Direction, req.Price, req.Quantity, req.TimeInForce)
		msgs := []sdk.Msg{msg}
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		bldr := sdkauth.NewTxBuilderFromCLI().
			WithTxEncoder(utils.GetTxEncoder(cdc)).
			WithKeybase(kb)

		bldr, sdkErr := utils.PrepareTxBuilder(bldr, ctx)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}

		broadcastResB, sdkErr := bldr.BuildAndSign(kb.GetName(), auth.MustGetKBPassphraseFromSession(r), msgs)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}
		broadcastRes, sdkErr := ctx.BroadcastTxCommit(broadcastResB)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}

		var orderIDStr string
		for _, log := range broadcastRes.Logs {
			if strings.HasPrefix(log.Log, "order_id") {
				orderIDStr = strings.TrimPrefix(log.Log, "order_id:")
				break
			}
		}
		orderID := sdk.NewUintFromString(orderIDStr)
		res := OrderCreationResponse{
			BlockInclusion: embedded.BlockInclusion{
				BlockNumber:     broadcastRes.Height,
				TransactionHash: broadcastRes.TxHash,
				BlockTimestamp:  broadcastRes.Timestamp,
			},
			ID:          orderID,
			MarketID:    msg.MarketID,
			Direction:   msg.Direction,
			Price:       msg.Price,
			Quantity:    msg.Quantity,
			Type:        req.Type,
			TimeInForce: msg.TimeInForce,
			Status:      "OPEN",
		}

		out, sdkErr := cdc.MarshalJSON(res)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}
		if _, err := w.Write(out); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
	}
}
