package store

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	headKey  = "head"
	valueKey = "value"
)

var ErrNoEntities = errors.New("no entities defined yet")

type Identifiable interface {
	GetID() EntityID
	SetID(id EntityID)
}

type Incrementing struct {
	backend KVStore
	cdc     *codec.Codec
}

func NewIncrementing(backend KVStore, cdc *codec.Codec) *Incrementing {
	return &Incrementing{
		backend: backend,
		cdc:     cdc,
	}
}

func (inc *Incrementing) ByID(id EntityID, val interface{}) error {
	b := inc.backend.Get(inc.ValueKey(id))
	if b == nil {
		return errors.New("not found")
	}
	return inc.cdc.UnmarshalBinaryBare(b, val)
}

func (inc *Incrementing) Head(val interface{}) error {
	head := inc.HeadID()
	if !head.IsDefined() {
		return ErrNoEntities
	}

	return inc.ByID(head, val)
}

func (inc *Incrementing) HasID(id EntityID) bool {
	return inc.backend.Has(inc.ValueKey(id))
}

func (inc *Incrementing) Iterator() sdk.Iterator {
	return KVStorePrefixIterator(inc.backend, []byte(valueKey))
}

func (inc *Incrementing) ReverseIterator() sdk.Iterator {
	return KVStoreReversePrefixIterator(inc.backend, []byte(valueKey))
}

func (inc *Incrementing) Insert(val Identifiable) error {
	if !val.GetID().IsZero() {
		return errors.New("id must be zero")
	}

	id := inc.HeadID().Inc()
	val.SetID(id)
	b, err := inc.cdc.MarshalBinaryBare(val)
	if err != nil {
		return err
	}
	inc.backend.Set(inc.ValueKey(id), b)
	inc.backend.Set(inc.HeadKey(), id.Bytes())
	return nil
}

func (inc *Incrementing) HeadID() EntityID {
	b := inc.backend.Get(inc.HeadKey())
	if b == nil {
		return NewEntityID(0)
	}

	return NewEntityIDFromBytes(b)
}

func (inc *Incrementing) HeadKey() []byte {
	return []byte(headKey)
}

func (inc *Incrementing) ValueKey(id EntityID) []byte {
	return PrefixKeyString(valueKey, id.Bytes())
}
