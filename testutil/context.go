package testutil

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DummyContext() sdk.Context {
	return sdk.NewContext(
		store.NewCommitMultiStore(db.NewMemDB()),
		abci.Header{ChainID: "unit-test-chain", Height: 1, Time: time.Unix(1558332092, 0)},
		false,
		log.NewNopLogger(),
	)
}
