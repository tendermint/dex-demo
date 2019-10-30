package fill

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/types"

	"github.com/tendermint/dex-demo/embedded"
	"github.com/tendermint/dex-demo/embedded/auth"
	"github.com/tendermint/dex-demo/pkg/conv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/user/fills", auth.DefaultAuthMW(userFills(ctx, cdc))).Methods("GET")
}

func userFills(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := auth.MustGetKBFromSession(r).GetAddr()
		q := r.URL.Query()

		var startBlock int
		var endBlock int
		var err error
		if start, ok := q["start_block"]; ok {
			startBlock, err = strconv.Atoi(start[0])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid start block")
				return
			}
		}
		if end, ok := q["end_block"]; ok {
			endBlock, err = strconv.Atoi(end[0])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid end block")
				return
			}
		}

		req := QueryRequest{
			Owner:      owner,
			StartBlock: int64(startBlock),
			EndBlock:   int64(endBlock),
		}
		fillsB, _, err := ctx.QueryWithData("custom/fill/get", cdc.MustMarshalBinaryBare(req))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var fills QueryResult
		cdc.MustUnmarshalJSON(fillsB, &fills)
		res := RESTQueryResult{
			Fills: make([]RESTFill, 0),
		}

		seenBlocks := make(map[int64]types.Block)
		for _, fill := range fills.Fills {
			var block types.Block
			var ok bool
			block, ok = seenBlocks[fill.BlockNumber]
			if !ok {
				b, err := getBlock(ctx, fill.BlockNumber)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				block = *b
				seenBlocks[fill.BlockNumber] = block
			}

			res.Fills = append(res.Fills, RESTFill{
				BlockInclusion: embedded.BlockInclusion{
					BlockNumber:    block.Height,
					BlockTimestamp: conv.FormatISO8601(block.Time),
				},
				QuantityFilled:   fill.QtyFilled,
				QuantityUnfilled: fill.QtyUnfilled,
				Direction:        fill.Direction,
				OrderID:          fill.OrderID,
				Pair:             fill.Pair,
				Price:            fill.Price,
				Owner:            fill.Owner,
			})
		}

		resB := cdc.MustMarshalJSON(res)
		embedded.PostProcessResponse(w, ctx, resB)
	}
}

func getBlock(ctx context.CLIContext, height int64) (*types.Block, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Block(&height)
	if err != nil {
		return nil, err
	}
	return res.Block, nil
}
