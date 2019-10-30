package session

import (
	"errors"
	"fmt"
	"net/http"
)

func MustGetStr(r *http.Request, key string) string {
	out, err := GetStr(r, key)
	if err != nil {
		panic(err)
	}
	return out
}

func GetStr(r *http.Request, key string) (string, error) {
	store, _ := SessionStore.Get(r, sessionName)
	val, ok := store.Values[key]
	if !ok || val == "" {
		return "", errors.New(fmt.Sprintf("key %s not found in session", key))
	}
	return val.(string), nil
}

func SetStrings(w http.ResponseWriter, r *http.Request, kvPairs ...string) error {
	if len(kvPairs) < 2 || len(kvPairs)%2 != 0 {
		return errors.New("mismatched KV pairs")
	}

	store, _ := SessionStore.Get(r, sessionName)
	for i := 0; i < len(kvPairs); i += 2 {
		k := kvPairs[i]
		v := kvPairs[i+1]
		store.Values[k] = v
	}

	return store.Save(r, w)
}
