package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseISO8601(t *testing.T) {
	in := "2019-07-08T00:50:29.000Z"
	_, err := ParseISO8601(in)
	assert.NoError(t, err)
}
