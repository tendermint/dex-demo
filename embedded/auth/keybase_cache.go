package auth

import (
	"errors"
	"net/http"
	"sync"

	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/dex-demo/embedded/session"
)

var kb *Keybase
var currID string
var mtx sync.RWMutex

func GetKBFromSession(r *http.Request) (*Keybase, error) {
	id, err := session.GetStr(r, keybaseIDKey)
	if err != nil {
		return nil, err
	}
	kb := GetKB(id)
	if kb == nil {
		return nil, errors.New("no keybase found")
	}
	return kb, nil
}

func MustGetKBFromSession(r *http.Request) *Keybase {
	kb, err := GetKBFromSession(r)
	if err != nil {
		panic(err)
	}
	return kb
}

func MustGetKBPassphraseFromSession(r *http.Request) string {
	return session.MustGetStr(r, keybasePassphraseKey)
}

func GetKB(id string) *Keybase {
	mtx.RLock()
	defer mtx.RUnlock()
	if currID != id {
		return nil
	}

	return kb
}

func ReplaceKB(name string, passphrase string, pk crypto.PrivKey) string {
	mtx.Lock()
	defer mtx.Unlock()
	currID = ReadStr32()
	kb = NewHotKeybase(name, passphrase, pk)
	return currID
}
