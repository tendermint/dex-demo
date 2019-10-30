package session

import (
	"github.com/gorilla/sessions"
)

const sessionName = "uex_session"

var SessionStore = sessions.NewCookieStore(generateSessionKey())

// TODO: pull from config
func generateSessionKey() []byte {
	var out [32]byte
	return out[:]
}
