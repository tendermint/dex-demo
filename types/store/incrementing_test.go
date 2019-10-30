package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/dex-demo/testutil/testflags"

	"github.com/cosmos/cosmos-sdk/codec"
)

type incrementingTest struct {
	ID  EntityID
	Foo string
	Bar int
}

func (it *incrementingTest) GetID() EntityID {
	return it.ID
}

func (it *incrementingTest) SetID(id EntityID) {
	it.ID = id
}

func TestIncrementing(t *testing.T) {
	testflags.UnitTest(t)
	db := dbm.NewMemDB()
	inc := NewIncrementing(db, codec.New())
	data := incrementingTest{
		ID:  NewEntityID(1),
		Foo: "hello",
		Bar: 1,
	}

	err := inc.Insert(&data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be zero")

	data.ID = NewEntityID(0)
	err = inc.Insert(&data)
	assert.NoError(t, err)

	var retrieved incrementingTest
	err = inc.ByID(data.ID, &retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "hello", retrieved.Foo)
	assert.Equal(t, 1, retrieved.Bar)
	assert.True(t, NewEntityID(1).Equals(retrieved.ID))
	assert.True(t, inc.HasID(retrieved.ID))

	err = inc.ByID(NewEntityID(999), &retrieved)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	data.ID = NewEntityID(0)
	err = inc.Insert(&data)
	assert.NoError(t, err)
	expID := NewEntityID(2)
	assert.True(t, inc.HasID(expID))
	assert.True(t, expID.Equals(inc.HeadID()))
}
