package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/dex-demo/testutil/testflags"

	"github.com/cosmos/cosmos-sdk/types"
)

type lastCall struct {
	key   []byte
	start []byte
	end   []byte
}

type dumbKVStore struct {
	last lastCall
}

func (d *dumbKVStore) Get(key []byte) []byte {
	d.last = lastCall{key: key}
	return nil
}

func (d *dumbKVStore) Has(key []byte) bool {
	d.last = lastCall{key: key}
	return true
}

func (d *dumbKVStore) Set(key, value []byte) {
	d.last = lastCall{key: key}
}

func (d *dumbKVStore) Delete(key []byte) {
	d.last = lastCall{key: key}
}

func (d *dumbKVStore) Iterator(start, end []byte) types.Iterator {
	d.last = lastCall{start: start, end: end}
	return nil
}

func (d *dumbKVStore) ReverseIterator(start, end []byte) types.Iterator {
	d.last = lastCall{start: start, end: end}
	return nil
}

func TestPrefixed(t *testing.T) {
	testflags.UnitTest(t)
	kvs := &dumbKVStore{}
	pref := NewPrefixed(kvs, []byte{0x44})

	pref.Get([]byte{0x01})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x01}, kvs.last.key))

	pref.Has([]byte{0x02})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x02}, kvs.last.key))

	pref.Set([]byte{0x03}, nil)
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x03}, kvs.last.key))

	pref.Delete([]byte{0x04})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x04}, kvs.last.key))

	pref.Iterator([]byte{0x05}, []byte{0x06})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x05}, kvs.last.start))
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x06}, kvs.last.end))

	pref.ReverseIterator([]byte{0x07}, []byte{0x08})
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x07}, kvs.last.start))
	assert.True(t, bytes.Equal([]byte{0x44, 0x2F, 0x08}, kvs.last.end))
}
