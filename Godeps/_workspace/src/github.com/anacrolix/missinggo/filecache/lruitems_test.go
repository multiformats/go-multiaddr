package filecache

import (
	"math/rand"
	"testing"
	"time"

	"github.com/jbenet/go-multiaddr-net/Godeps/_workspace/src/github.com/bradfitz/iter"
)

func BenchmarkInsert(b *testing.B) {
	for _ = range iter.N(b.N) {
		li := newLRUItems()
		for _ = range iter.N(10000) {
			r := rand.Int63()
			t := time.Unix(r/1e9, r%1e9)
			li.Insert(ItemInfo{
				Accessed: t,
			})
		}
	}
}
