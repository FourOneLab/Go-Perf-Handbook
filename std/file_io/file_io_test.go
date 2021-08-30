package file_io

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func BenchmarkWriteFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Create("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		for i := 0; i < 100000; i++ {
			f.WriteString("some text!\n")
		}
		f.Close()
	}
}

func BenchmarkWriteFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Create("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		w := bufio.NewWriter(f)
		for i := 0; i < 100000; i++ {
			w.WriteString("some text!\n")
		}
		w.Flush()
		f.Close()
	}
}

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("/tmp/test.txt")
		if err != nil {
			panic(err)
		}
		b := make([]byte, 10)
		for err == nil {
			_, err = f.Read(b)
		}
		if err != io.EOF {
			panic(err)
		}
		f.Close()
	}
}

func BenchmarkReadFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("/tmp/test.txt")
		if err != nil {
			panic(err)
		}

		r := bufio.NewReader(f)
		for err == nil {
			_, err = r.ReadString('\n')
		}
		if err != io.EOF {
			panic(err)
		}
		f.Close()
	}
}
