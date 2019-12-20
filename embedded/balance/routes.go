package balance

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/tendermint/dex-demo/embedded"
	"github.com/tendermint/dex-demo/embedded/auth"
	"github.com/tendermint/dex-demo/x/asset/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authsdk "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec, enableFaucet bool) {
	r.Handle("/user/balances", auth.DefaultAuthMW(getBalanceHandler(ctx, cdc))).Methods("GET")
	r.Handle("/user/transfer", auth.LoginRequiredMW(auth.OTPRequiredMW(transferBalanceHandler(ctx, cdc)))).Methods("POST")

	if enableFaucet {
		r.Handle("/faucet/transfer", faucetHandler(ctx, cdc)).Methods("POST")
	}
}

func getBalanceHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := auth.MustGetKBFromSession(r)

		req := GetQueryRequest{
			Address: owner.GetAddr(),
		}

		resB, _, err := ctx.QueryWithData("custom/balance/get", cdc.MustMarshalBinaryBare(req))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		embedded.PostProcessResponse(w, ctx, resB)
	}
}

func transferBalanceHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferBalanceRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		kb := auth.MustGetKBFromSession(r)
		doTransfer(kb, ctx, w, cdc, req.To, req.Amount, req.AssetID, auth.MustGetKBPassphraseFromSession(r))
	}
}

func faucetHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferBalanceRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Auth header must be provided.", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, ":")
		if len(parts) != 2 {
			http.Error(w, "Auth header must be formatted as username:password.", http.StatusUnauthorized)
			return
		}

		username := parts[0]
		passphrase := parts[1]

		if username != auth.AccountName {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		diskKB, err := keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pk, err := diskKB.ExportPrivateKeyObject(auth.AccountName, passphrase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		kb := auth.NewHotKeybase(auth.AccountName, passphrase, pk)
		doTransfer(kb, ctx, w, cdc, req.To, req.Amount, req.AssetID, passphrase)
	}
}

func doTransfer(kb *auth.Keybase, ctx context.CLIContext, w http.ResponseWriter, cdc *codec.Codec, to sdk.AccAddress, amount sdk.Uint, assetID sdk.Uint, passphrase string) {
	owner := kb.GetAddr()
	ctx = ctx.WithFromAddress(owner)
	msg := types.NewMsgTransfer(assetID, owner, to, amount)
	msgs := []sdk.Msg{msg}
	err := msg.ValidateBasic()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	bldr := authsdk.NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithKeybase(kb)
	bldr, sdkErr := utils.PrepareTxBuilder(bldr, ctx)
	if sdkErr != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
		return
	}
	broadcastResB, sdkErr := bldr.BuildAndSign(kb.GetName(), passphrase, msgs)
	if sdkErr != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
		return
	}
	broadcastRes, sdkErr := ctx.BroadcastTxCommit(broadcastResB)
	if sdkErr != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
		return
	}
	res := TransferBalanceResponse{
		BlockInclusion: embedded.BlockInclusion{
			BlockNumber:     broadcastRes.Height,
			TransactionHash: broadcastRes.TxHash,
			BlockTimestamp:  broadcastRes.Timestamp,
		},
	}
	out, sdkErr := cdc.MarshalJSON(res)
	if sdkErr != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
	}
	if _, err := w.Write(out); err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
