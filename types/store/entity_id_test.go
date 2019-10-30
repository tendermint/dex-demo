package store

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityID(t *testing.T) {
	t.Run("should be stringifiable", func(t *testing.T) {
		assert.Equal(t, "10", NewEntityID(10).String())
	})
	t.Run("should be instantiable from a string", func(t *testing.T) {
		assert.EqualValues(t, NewEntityID(1), NewEntityIDFromString("1"))
	})
	t.Run("should be incrementable without mutation", func(t *testing.T) {
		a := NewEntityID(10)
		assert.EqualValues(t, NewEntityID(11), a.Inc())
		assert.EqualValues(t, NewEntityID(10), a)
	})
	t.Run("should return IsDefined()", func(t *testing.T) {
		assert.False(t, NewEntityID(0).IsDefined())
		assert.True(t, NewEntityID(1).IsDefined())
	})
	t.Run("should define equality", func(t *testing.T) {
		assert.True(t, NewEntityID(1).Equals(NewEntityID(1)))
		assert.False(t, NewEntityID(2).Equals(NewEntityID(1)))
	})
	t.Run("should return a fixed length bytes representation", func(t *testing.T) {
		res := hex.EncodeToString(NewEntityID(1000).Bytes())
		assert.Equal(t, 64, len(res))
		assert.EqualValues(t, "00000000000000000000000000000000000000000000000000000000000003e8", res)
	})
}
