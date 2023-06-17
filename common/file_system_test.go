package common

import (
	"testing"
)

func TestOpenFile(t *testing.T) {
	fs := &FileSystem{}
	fs.OpenFile("a.txt", WriteOnly|Create, NoLock)
}

func TestRead(t *testing.T) {
	fs := &FileSystem{}
	handle := fs.OpenFile("a.txt", ReadOnly, NoLock)
	buffer := make([]byte, 0, 4096)
	handle.Read(&buffer, 2, 0)
	expect := []byte{'H', 'i'}

	for i := range expect {
		if expect[i] != buffer[i] {
			t.Errorf("Expect %c and got %c", expect[i], buffer[i])
		}
	}
}

func TestWrite(t *testing.T) {
	fs := &FileSystem{}
	text := "hello world"
	handle := fs.OpenFile("b.txt", WriteOnly|Create, NoLock)
	buffer := []byte(text)
	handle.Write(buffer, uint64(len(text)), 0)

	buffer = make([]byte, 0, 4096)
	handle.Read(&buffer, uint64(len(text)), 0)

	for i := range text {
		if text[i] != buffer[i] {
			t.Errorf("Expect %c, got %c\n", text[i], buffer[i])
		}
	}
}
