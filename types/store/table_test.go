package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tm-db"
)

func TestTable(t *testing.T) {
	db := dbm.NewMemDB()
	tb1 := NewTable(db, "foo")

	k1 := []byte{0x00}
	t.Run("should be able to get and set", func(t *testing.T) {
		tb1.Set(k1, []byte("value"))
		assert.Equal(t, "value", string(tb1.Get(k1)))
	})
	t.Run("set should overwrite", func(t *testing.T) {
		tb1.Set(k1, []byte("value1"))
		tb1.Set(k1, []byte("value2"))
		assert.Equal(t, "value2", string(tb1.Get(k1)))
	})
	t.Run("keys should be namespaced", func(t *testing.T) {
		tb2 := NewTable(db, "bar")
		tb1.Set(k1, []byte("value1"))
		tb2.Set(k1, []byte("value2"))
		assert.Equal(t, "value1", string(tb1.Get(k1)))
		assert.Equal(t, "value2", string(tb2.Get(k1)))
	})
	t.Run("should be able to delete", func(t *testing.T) {
		tb1.Set(k1, []byte("value"))
		tb1.Delete(k1)
		assert.Nil(t, tb1.Get(k1))
		assert.False(t, tb1.Has(k1))
	})
	t.Run("has returns true for existing keys and false for non existent keys", func(t *testing.T) {
		tb1.Set(k1, []byte("value"))
		k2 := []byte{0x01}
		assert.True(t, tb1.Has(k1))
		assert.False(t, tb1.Has(k2))
	})
	t.Run("can iterate between values and stop part way", func(t *testing.T) {
		for i := 0; i <= 255; i++ {
			tb1.Set([]byte{byte(i)}, []byte{byte(i)})
		}

		// note: end is exclusive
		last := 0
		tb1.Iterator([]byte{0x01}, []byte{0x23}, func(k []byte, v []byte) bool {
			i := last + 1
			assert.Equal(t, 1, len(k))
			assert.Equal(t, 1, len(v))
			assert.Equal(t, i, int(k[0]))
			assert.Equal(t, i, int(v[0]))
			last = int(v[0])
			return true
		})
		assert.Equal(t, 34, last)

		count := 0
		tb1.Iterator([]byte{0x01}, []byte{0x23}, func(k []byte, v []byte) bool {
			count++
			return count < 32
		})
		assert.Equal(t, 32, count)
	})
	t.Run("can reverse iterate between values and stop part way", func(t *testing.T) {
		for i := 0; i <= 255; i++ {
			tb1.Set([]byte{byte(i)}, []byte{byte(i)})
		}

		// end is exclusive.
		last := 35
		tb1.ReverseIterator([]byte{0x01}, []byte{0x23}, func(k []byte, v []byte) bool {
			i := last - 1
			assert.Equal(t, 1, len(k))
			assert.Equal(t, 1, len(v))
			assert.Equal(t, i, int(k[0]))
			assert.Equal(t, i, int(v[0]))
			last = int(v[0])
			return true
		})
		assert.Equal(t, 1, last)

		count := 0
		tb1.ReverseIterator([]byte{0x01}, []byte{0x23}, func(k []byte, v []byte) bool {
			count++
			return count < 7
		})
		assert.Equal(t, 7, count)
	})
	t.Run("can iterate over a prefix", func(t *testing.T) {
		for i := 0; i < 255; i++ {
			tb1.Set(PrefixKeyString("pref1", []byte{byte(i)}), []byte{byte(i)})
			tb1.Set(PrefixKeyString("pref2", []byte{byte(i)}), []byte{byte(i)})
		}

		i := 0
		tb1.PrefixIterator([]byte("pref1"), func(k []byte, v []byte) bool {
			expK := append([]byte("pref1/"), byte(i))
			assert.Equal(t, 7, len(k))
			assert.Equal(t, 1, len(v))
			assert.EqualValues(t, expK, k)
			assert.EqualValues(t, i, int(byte(v[0])))
			i++
			return true
		})
		assert.Equal(t, 255, i)
	})
	t.Run("can reverse iterate over a prefix", func(t *testing.T) {
		for i := 0; i < 255; i++ {
			tb1.Set(PrefixKeyString("pref1", []byte{byte(i)}), []byte{byte(i)})
			tb1.Set(PrefixKeyString("pref2", []byte{byte(i)}), []byte{byte(i)})
		}

		i := 254
		tb1.ReversePrefixIterator([]byte("pref1"), func(k []byte, v []byte) bool {
			expK := append([]byte("pref1/"), byte(i))
			assert.Equal(t, 7, len(k))
			assert.Equal(t, 1, len(v))
			assert.Equal(t, expK, k)
			assert.EqualValues(t, i, int(byte(v[0])))
			i--
			return true
		})
		// there are 255 entities, so this goes to -1
		assert.Equal(t, -1, i)
	})
}
