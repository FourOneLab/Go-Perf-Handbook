package templates

import (
	"io"
	"testing"
	"text/template"
)

var bookTemplate = `
Title: {{ .Title }},
Author: {{ .Author }}
{{ if .Pages }}
Number of pages: {{ .Pages }}.
{{ end }}
{{ range .Chapters }}
{{ . }},
{{ end }}
`

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

func BenchmarkExecute(b *testing.B) {
	book := &Book{
		Title:    "The Art of Computer Programming, Vol. 3",
		Author:   "Donald E. Knuth",
		Pages:    800,
		Chapters: []string{"Sorting", "Searching"},
	}

	t := template.Must(template.New("book").Parse(bookTemplate))

	for i := 0; i < b.N; i++ {
		if err := t.Execute(io.Discard, book); err != nil {
			panic(err)
		}
	}
}
