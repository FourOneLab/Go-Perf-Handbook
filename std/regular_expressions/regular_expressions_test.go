package regular_expressions

import (
	"regexp"
	"testing"
)

var testRegexp = `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]+$`

func BenchmarkMatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := regexp.MatchString(testRegexp, "qwer@example.com"); err != nil {
			panic(err)
		}
	}
}

func BenchmarkMatchStringCompiled(b *testing.B) {
	compile, err := regexp.Compile(testRegexp)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compile.MatchString("qwer@example.com")
	}
}
