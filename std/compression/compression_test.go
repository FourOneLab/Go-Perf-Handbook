package compression

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	data, err := os.ReadFile("/tmp/test.txt")
	if err != nil {
		panic(err)
	}
	zipWriter := gzip.NewWriter(io.Discard)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err = zipWriter.Write(data); err != nil {
			panic(err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	data, err := os.ReadFile("/tmp/test.txt")
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	zipWriter := gzip.NewWriter(buf)
	if _, err = zipWriter.Write(data); err != nil {
		panic(err)
	}
	if err = zipWriter.Close(); err != nil {
		panic(err)
	}

	reader := bytes.NewReader(buf.Bytes())
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Reset(buf.Bytes())
		gzipReader.Reset(reader)
		if _, err = ioutil.ReadAll(gzipReader); err != nil {
			panic(err)
		}
	}
}
