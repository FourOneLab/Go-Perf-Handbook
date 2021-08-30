package random_numbers

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"
)

func BenchmarkMathRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int63()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := crand.Int(crand.Reader, big.NewInt(27)); err != nil {
			panic(err)
		}
	}
}
