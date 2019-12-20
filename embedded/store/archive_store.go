package store

type IteratorCB func(k []byte, v []byte) bool

type ArchiveStore interface {
	Get(key []byte) []byte
	Has(key []byte) bool
	Set(key []byte, value []byte)
	Delete(key []byte)
	Iterator(start []byte, end []byte, cb IteratorCB)
	ReverseIterator(start []byte, end []byte, cb IteratorCB)
	PrefixIterator(start []byte, cb IteratorCB)
	ReversePrefixIterator(start []byte, cb IteratorCB)
	Substore(prefix string) ArchiveStore
}
