package store

import (
	"bytes"
	"encoding/binary"
)

func PrefixKeyString(prefix string, subkeys ...[]byte) []byte {
	buf := [][]byte{[]byte(prefix)}
	return PrefixKeyBytes(append(buf, subkeys...)...)
}

func PrefixKeyBytes(subkeys ...[]byte) []byte {
	if len(subkeys) == 0 {
		return []byte{}
	}

	var buf bytes.Buffer
	buf.Write(subkeys[0])

	if len(subkeys) > 1 {
		for _, sk := range subkeys[1:] {
			if len(sk) == 0 {
				continue
			}

			buf.WriteRune('/')
			buf.Write(sk)
		}
	}

	return buf.Bytes()
}

func IntSubkey(subkey int) []byte {
	if subkey < 0 {
		panic("cannot use negative numbers in subkeys")
	}
	return Uint64Subkey(uint64(subkey))
}

func Int64Subkey(subkey int64) []byte {
	if subkey < 0 {
		panic("cannot use negative numbers in subkeys")
	}
	return Uint64Subkey(uint64(subkey))
}

func Uint64Subkey(subkey uint64) []byte {
	b := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, subkey)
	return b
}
