package numeric_conversions

import (
	"strconv"
	"testing"
)

func BenchmarkParseBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := strconv.ParseBool("true"); err != nil {
			panic(err)
		}
	}
}

func BenchmarkParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := strconv.ParseInt("7182818284", 10, strconv.IntSize); err != nil {
			panic(err)
		}
	}
}

func BenchmarkParseFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := strconv.ParseFloat("3.141592653", strconv.IntSize); err != nil {
			panic(err)
		}
	}
}
