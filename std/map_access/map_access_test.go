package map_access

import (
	"math/rand"
	"strconv"
	"testing"
)

const numItems = 1000000

func BenchmarkMapStringKeys(b *testing.B) {
	m := make(map[string]string)
	k := make([]string, 0)

	for i := 0; i < numItems; i++ {
		key := strconv.Itoa(rand.Intn(numItems))
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	n := 0
	l := len(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m[k[n]]; ok {
		}
		n++
		if n >= l {
			n = 0
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	m := make(map[int]string)
	k := make([]int, 0)

	for i := 0; i < numItems; i++ {
		key := rand.Intn(numItems)
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	n := 0
	l := len(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m[k[n]]; ok {
		}
		n++
		if n >= l {
			n = 0
		}
	}
}
