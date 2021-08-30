package sorting

import (
	"math/rand"
	"sort"
	"testing"
)

func generateSlice(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, rand.Intn(1e9))
	}
	return s
}

func BenchmarkSort1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := generateSlice(1000)
		b.StartTimer()
		sort.Ints(s)
	}
}

func BenchmarkSort100000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := generateSlice(100000)
		b.StartTimer()
		sort.Ints(s)
	}
}

func BenchmarkSort10000000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := generateSlice(10000000)
		b.StartTimer()
		sort.Ints(s)
	}
}
