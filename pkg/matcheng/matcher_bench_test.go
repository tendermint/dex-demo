package matcheng

import (
	"math/rand"
	"testing"

	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BenchmarkMatching(b *testing.B) {
	id := store.NewEntityID(0)
	matcher := GetMatcher()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		matcher.Reset()
		for j := 0; j < 10000; j++ {
			id = id.Inc()
			matcher.EnqueueOrder(Bid, id, sdk.NewUint(uint64(j)), sdk.NewUint(uint64(j)))
		}
		for j := 100; j < 11000; j++ {
			id := id.Inc()
			matcher.EnqueueOrder(Ask, id, sdk.NewUint(uint64(j)), sdk.NewUint(uint64(j)))
		}
		b.StartTimer()
		matcher.Match()
	}
}

func BenchmarkQueueing(b *testing.B) {
	id := store.NewEntityID(0)
	matcher := GetMatcher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		matcher.Reset()
		b.StartTimer()
		for j := 0; j < 100; j++ {
			id = id.Inc()
			price := sdk.NewUint(rand.Uint64())
			quantity := sdk.NewUint(rand.Uint64())
			matcher.EnqueueOrder(Bid, id.Inc(), price, quantity)
		}
	}
}
