package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func ReadStr32() string {
	return ReadStrN(32)
}

func ReadStrN(byteLen int) string {
	return hex.EncodeToString(ReadN(byteLen))
}

func ReadN(n int) []byte {
	buf := make([]byte, n, n)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}
