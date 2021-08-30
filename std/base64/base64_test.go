package base64

import (
	"encoding/base64"
	"math/rand"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	data := make([]byte, 1024)
	rand.Read(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base64.StdEncoding.EncodeToString(data)
	}
}

func BenchmarkDecode(b *testing.B) {
	data := make([]byte, 1024)
	rand.Read(data)
	encoded := base64.StdEncoding.EncodeToString(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := base64.StdEncoding.DecodeString(encoded); err != nil {
			panic(err)
		}
	}
}
