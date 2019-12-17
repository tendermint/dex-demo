package auth

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/tendermint/dex-demo/embedded/session"
)

const (
	keybaseIDKey         = "keybaseID"
	keybasePassphraseKey = "keybasePassphrase"
	otpHeader            = "X-OTP-Token"
)

func DefaultAuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LoginRequiredMW(next).ServeHTTP(w, r)
	})
}

func LoginRequiredMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.SessionStore.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		kbID, ok := store.Values[keybaseIDKey]
		if !ok || GetKB(kbID.(string)) == nil {
			http.Error(w, "Not logged in.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OTPRequiredMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(otpHeader)
		if header == "" {
			http.Error(w, "No OTP header provided.", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func HandleCORSMW(next http.Handler) http.Handler {
	// TODO: Pull from config
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(next)
}
