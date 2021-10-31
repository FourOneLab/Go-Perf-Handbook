package serialization

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

func generateObject() *Book {
	return &Book{
		Title:    "The Art of Computer Programming, Vol. 2",
		Author:   "Donald E. Knuth",
		Pages:    784,
		Chapters: []string{"Random numbers", "Arithmetic"},
	}
}

func generatePBObject() *BookProto {
	book := generateObject()
	return &BookProto{
		Title:    *proto.String(book.Title),
		Author:   *proto.String(book.Author),
		Pages:    *proto.Int64(int64(book.Pages)),
		Chapters: book.Chapters,
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	book := generateObject()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(book)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	out, err := json.Marshal(generateObject())
	if err != nil {
		panic(err)
	}
	book := &Book{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(out, book); err != nil {
			panic(err)
		}
	}
}

func BenchmarkPBMarshal(b *testing.B) {
	bookProto := generatePBObject()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := proto.Marshal(bookProto); err != nil {
			panic(err)
		}
	}
}

func BenchmarkPBUnmarshal(b *testing.B) {
	out, err := proto.Marshal(generatePBObject())
	if err != nil {
		panic(err)
	}
	bookProto := &BookProto{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = proto.Unmarshal(out, bookProto); err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobMarshal(b *testing.B) {
	book := generateObject()
	encoder := gob.NewEncoder(ioutil.Discard)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := encoder.Encode(book); err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobUnmarshal(b *testing.B) {
	book := generateObject()
	buf := &bytes.Buffer{}
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(book); err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := decoder.Decode(&Book{}); err != nil {
			return
		}
	}
}
