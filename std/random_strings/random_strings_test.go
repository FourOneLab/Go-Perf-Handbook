package random_strings

import (
	crand "crypto/rand"
	"math/rand"
	"testing"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

func randomBytes(n int) []byte {
	byteSlice := make([]byte, n)
	if _, err := rand.Read(byteSlice); err != nil {
		panic(err)
	}
	return byteSlice
}

func cryptoRandomBytes(n int) []byte {
	byteSlice := make([]byte, n)
	if _, err := crand.Read(byteSlice); err != nil {
		panic(err)
	}
	return byteSlice
}

func randomString(byteSlice []byte) string {
	for i, b := range byteSlice {
		byteSlice[i] = letters[b%64]
	}
	return string(byteSlice)
}

func BenchmarkMatchRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomString(randomBytes(16))
	}
}

func BenchmarkCryptoRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomString(cryptoRandomBytes(16))
	}
}
