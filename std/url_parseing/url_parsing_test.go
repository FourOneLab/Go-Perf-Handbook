package url_parseing

import (
	"net/url"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	testURL := "https://www.example.com/path/file.html?param1=value1&param2=123"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := url.Parse(testURL); err != nil {
			panic(err)
		}
	}
}
