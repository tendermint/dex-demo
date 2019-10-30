package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/dex-demo/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/dex-demo/embedded/session"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	sub := r.PathPrefix("/auth").Subrouter()
	sub.HandleFunc("/login", loginHandler()).Methods("POST")
	sub.Handle("/logout", DefaultAuthMW(logoutHandler())).Methods("POST")
	sub.HandleFunc("/csrf_token", csrfTokenHandler()).Methods("GET")
	sub.Handle("/me", DefaultAuthMW(meHandler(ctx, cdc))).Methods("GET")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var req LoginRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Username != AccountName {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		kbID, hotPW, err := authorize(req.Password)
		if err != nil {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		err = session.SetStrings(w, r, keybaseIDKey, kbID, keybasePassphraseKey, hotPW)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store, err := session.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete(store.Values, keybaseIDKey)
		err = store.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func csrfTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := GetCSRFToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(tok))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type MeResponse struct {
	Address string `json:"address"`
}

func meHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := MustGetKBFromSession(r)
		addr := owner.GetAddr().String()
		res := &MeResponse{Address: addr}
		resB := cdc.MustMarshalJSON(res)
		embedded.PostProcessResponse(w, ctx, resB)
	}
}

func authorize(passphrase string) (string, string, error) {
	kb, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		return "", "", err
	}

	pk, err := kb.ExportPrivateKeyObject(AccountName, passphrase)
	if err != nil {
		return "", "", err
	}

	hotPassphrase := ReadStr32()
	return ReplaceKB(AccountName, hotPassphrase, pk), hotPassphrase, nil
}
