package string_concateation

import (
	"bytes"
	"strings"
	"testing"
)

var strLen = 1000

func BenchmarkConcatString(b *testing.B) {
	var str string
	n := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str += "x"
		n++
		if n >= strLen {
			n = 0
			str = ""
		}
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	var buffer bytes.Buffer
	n := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.WriteString("x")
		n++
		if n >= strLen {
			n = 0
			buffer = bytes.Buffer{}
		}
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	var builder strings.Builder
	n := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.WriteString("x")
		n++
		if n >= strLen {
			n = 0
			builder = strings.Builder{}
		}
	}
}
