package testutil

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TmpDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", strconv.Itoa(int(time.Now().UnixNano())))
	require.NoError(t, err)

	return dir, func() {
		err := os.RemoveAll(dir)
		require.NoError(t, err)
	}
}
