package log

import (
	"os"

	"github.com/tendermint/tendermint/libs/log"
)

func WithModule(module string) log.Logger {
	return log.NewTMLogger(os.Stdout).With("module", module)
}
