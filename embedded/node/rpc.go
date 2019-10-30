package node

import (
	core_types "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func LatestBlock(ctx context.CLIContext) (*core_types.ResultBlock, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}
	return node.Block(nil)
}
