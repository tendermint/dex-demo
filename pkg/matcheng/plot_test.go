package matcheng

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/dex-demo/testutil/testflags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPlotCurves(t *testing.T) {
	testflags.UnitTest(t)
	expected := `"Ask"
2 0
2 10
3 10
3 20
4 20
4 30


"Bid"
3 0
3 10
2 10
2 20
1 20
1 30
0 30
`

	res := &MatchResults{
		BidAggregates: []AggregatePrice{
			{sdk.OneUint(), sdk.NewUint(30)},
			{sdk.NewUint(2), sdk.NewUint(20)},
			{sdk.NewUint(3), sdk.NewUint(10)},
		},
		AskAggregates: []AggregatePrice{
			{sdk.NewUint(2), sdk.NewUint(10)},
			{sdk.NewUint(3), sdk.NewUint(20)},
			{sdk.NewUint(4), sdk.NewUint(30)},
		},
	}

	actual := PlotCurves(res.BidAggregates, res.AskAggregates)
	assert.Equal(t, expected, actual)

}
