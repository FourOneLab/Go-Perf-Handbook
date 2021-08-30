package object_creation

import (
	"sync"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

var pool = sync.Pool{New: func() interface{} { return &Book{} }}

func BenchmarkNoPool(b *testing.B) {
	book := new(Book)
	for i := 0; i < b.N; i++ {
		book = &Book{
			Title:  "The Art of Computer Programming, Vol. 1",
			Author: "Donald E. Knuth",
			Pages:  672,
		}
	}
	_ = book
}

func BenchmarkPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		book := pool.Get().(*Book)
		book.Title = "The Art of Computer Programming, Vol. 1"
		book.Author = "Donald E. Knuth"
		book.Pages = 672
		pool.Put(book)
	}
}
