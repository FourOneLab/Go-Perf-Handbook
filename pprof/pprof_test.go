package pprof

import (
	"math/rand"
	"testing"
)

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int63()
	}
}
