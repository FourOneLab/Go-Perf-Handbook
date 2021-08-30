package slice_appending

import "testing"

var numItems = 1000000

func BenchmarkSliceAppend(b *testing.B) {
	s := make([]byte, 0)
	n := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = append(s, 1)
		n++
		if n == numItems {
			b.StopTimer()
			n = 0
			s = make([]byte, 0)
			b.StartTimer()
		}
	}
}

func BenchmarkSliceAppendPreAlloc(b *testing.B) {
	s := make([]byte, 0, numItems)
	n := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = append(s, 1)
		n++
		if n == numItems {
			b.StopTimer()
			n = 0
			s = make([]byte, 0, numItems)
			b.StartTimer()
		}
	}
}
