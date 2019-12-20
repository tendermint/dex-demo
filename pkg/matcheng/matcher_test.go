package matcheng

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/testutil"
	"github.com/tendermint/dex-demo/testutil/testflags"
)

func TestMatcher_Golden(t *testing.T) {
	testflags.UnitTest(t)
	runTestFile(t)
}

func TestMatcher_NoClearingPrice(t *testing.T) {
	testflags.UnitTest(t)
	bids := [][2]uint64{
		{1, 10},
		{2, 10},
		{3, 10},
	}

	asks := [][2]uint64{
		{4, 10},
		{5, 10},
		{6, 10},
	}

	res, _ := doMatch(bids, asks)
	assert.Nil(t, res)
}

func TestMatcher_NoMarket(t *testing.T) {
	testflags.UnitTest(t)
	matcher := GetMatcher()
	defer ReturnMatcher(matcher)

	orders := [][2]uint64{
		{1, 10},
		{2, 10},
		{3, 10},
	}

	res, _ := doMatch(orders, nil)
	assert.Nil(t, res)
	matcher.Reset()
	res, _ = doMatch(nil, orders)
	assert.Nil(t, res)
}

func doMatch(bids [][2]uint64, asks [][2]uint64) (*MatchResults, map[string]Fill) {
	matcher := GetMatcher()
	defer ReturnMatcher(matcher)

	id := sdk.ZeroUint()
	if bids != nil {
		for _, bid := range bids {
			id = id.Add(sdk.OneUint())
			matcher.EnqueueOrder(Bid, id, sdk.NewUint(bid[0]), sdk.NewUint(bid[1]))
		}
	}
	if asks != nil {
		for _, ask := range asks {
			id = id.Add(sdk.OneUint())
			matcher.EnqueueOrder(Ask, id, sdk.NewUint(ask[0]), sdk.NewUint(ask[1]))
		}
	}

	res := matcher.Match()
	if res == nil {
		return nil, nil
	}

	fillsMap := make(map[string]Fill)
	for _, fill := range res.Fills {
		fillsMap[fill.OrderID.String()] = fill
	}

	return res, fillsMap
}

func assertFill(t *testing.T, fill Fill, qtyFilled uint64, qtyUnfilled uint64) {
	assert.NotNil(t, fill)
	testutil.AssertEqualUints(t, sdk.NewUint(qtyFilled), fill.QtyFilled, "order id %s has invalid qty filled", fill.OrderID)
	testutil.AssertEqualUints(t, sdk.NewUint(qtyUnfilled), fill.QtyUnfilled, "order id %s has invalid qty unfilled", fill.OrderID)
}

type testCase struct {
	name         string
	bids         [][2]uint64
	asks         [][2]uint64
	expectations [][3]uint64
	clearing     uint64
	lineNum      int
}

func runTestFile(t *testing.T) {
	f, err := os.Open(path.Join("testdata/matches.txt"))
	require.NoError(t, err)
	s := bufio.NewScanner(f)

	var currCase *testCase
	var cases []*testCase

	var lineNum int
	for s.Scan() {
		lineNum++
		line := s.Text()
		if line == "" {
			continue
		}
		firstChar := string(line[0])

		if firstChar == "#" {
			currCase = &testCase{
				name:    strings.Replace(line, "# ", "", 1),
				lineNum: lineNum,
			}
			cases = append(cases, currCase)
			continue
		}

		if currCase == nil {
			panic("at least one test case is required")
		}

		switch firstChar {
		case "b":
			split := strings.Split(line, " ")
			currCase.bids = append(currCase.bids, [2]uint64{mustParseUint(split[1]), mustParseUint(split[2])})
		case "a":
			split := strings.Split(line, " ")
			currCase.asks = append(currCase.asks, [2]uint64{mustParseUint(split[1]), mustParseUint(split[2])})
		case "c":
			split := strings.Split(line, " ")
			if currCase.clearing != 0 {
				panic("already set clearing for this test case")
			}
			currCase.clearing = mustParseUint(split[1])
		case "e":
			split := strings.Split(line, " ")
			currCase.expectations = append(currCase.expectations, [3]uint64{mustParseUint(split[1]), mustParseUint(split[2]), mustParseUint(split[3])})
		case "\n":
			// noop
		default:
			panic("unknown line")
		}
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			res, fillsMap := doMatch(tCase.bids, tCase.asks)
			assert.NotNil(t, res)
			assert.Equal(t, len(tCase.expectations), len(res.Fills), "number of fills do not match. line num: %d", tCase.lineNum)
			for _, exp := range tCase.expectations {
				assertFill(t, fillsMap[strconv.FormatUint(exp[0], 10)], exp[1], exp[2])
			}
		})
	}
}

func mustParseUint(s string) uint64 {
	out, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return out
}
