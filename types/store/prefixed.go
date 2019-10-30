package store

import sdk "github.com/cosmos/cosmos-sdk/types"

type Prefixed struct {
	backend KVStore
	prefix  []byte
}

func NewPrefixed(backend KVStore, prefix []byte) *Prefixed {
	return &Prefixed{
		backend: backend,
		prefix:  prefix,
	}
}

func (p *Prefixed) Get(key []byte) []byte {
	return p.backend.Get(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Has(key []byte) bool {
	return p.backend.Has(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Set(key, value []byte) {
	p.backend.Set(PrefixKeyBytes(p.prefix, key), value)
}

func (p *Prefixed) Delete(key []byte) {
	p.backend.Delete(PrefixKeyBytes(p.prefix, key))
}

func (p *Prefixed) Iterator(start, end []byte) sdk.Iterator {
	return p.backend.Iterator(PrefixKeyBytes(p.prefix, start), PrefixKeyBytes(p.prefix, end))
}

func (p *Prefixed) ReverseIterator(start, end []byte) sdk.Iterator {
	return p.backend.ReverseIterator(PrefixKeyBytes(p.prefix, start), PrefixKeyBytes(p.prefix, end))
}
