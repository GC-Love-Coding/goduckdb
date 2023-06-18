package common

import "testing"

// func TestOpenFile(t *testing.T) {
// 	fs := &FileSystem{}
// 	fs.OpenFile("a.txt", WriteOnly|Create, NoLock)
// }

// func TestRead(t *testing.T) {
// 	fs := &FileSystem{}
// 	handle := fs.OpenFile("a.txt", ReadOnly, NoLock)
// 	buffer := make([]byte, 0, 4096)
// 	handle.Read(&buffer, 2, 0)
// 	expect := []byte{'H', 'i'}

// 	for i := range expect {
// 		if expect[i] != buffer[i] {
// 			t.Errorf("Expect %c and got %c", expect[i], buffer[i])
// 		}
// 	}
// }

// func TestWrite(t *testing.T) {
// 	fs := &FileSystem{}
// 	text := "hello world"
// 	handle := fs.OpenFile("b.txt", WriteOnly|Create, NoLock)
// 	buffer := []byte(text)
// 	handle.Write(buffer, uint64(len(text)), 0)

// 	buffer = make([]byte, 0, 4096)
// 	handle.Read(&buffer, uint64(len(text)), 0)

// 	for i := range text {
// 		if text[i] != buffer[i] {
// 			t.Errorf("Expect %c, got %c\n", text[i], buffer[i])
// 		}
// 	}
// }

// func TestGetFileSize(t *testing.T) {
// 	fs := &FileSystem{}
// 	handle := fs.OpenFile("c.txt", WriteOnly|Create, NoLock)
// 	size := fs.GetFileSize(handle)

// 	if size != 0 {
// 		t.Errorf("Expect file size: 0, got: %d\n", size)
// 	}
// }

// func TestDirectoryExists(t *testing.T) {
// 	fs := &FileSystem{}
// 	directory := "/tmp/@dir@"

// 	if fs.DirectoryExists(directory) {
// 		t.Errorf("Expect directory %s not exist\n", directory)
// 	}
// }

// func TestRemoveDirectory(t *testing.T) {
// 	fs := &FileSystem{}
// 	directory := "/tmp/@dir@"
// 	fs.RemoveDirectory(directory)

// 	directory = "bin"
// 	fs.CreateDirectory(directory)
// 	fs.RemoveDirectory(directory)
// }

// func TestMoveFile(t *testing.T) {
// 	fs := &FileSystem{}
// 	fs.MoveFile("a.txt", "d.txt")
// }

func TestFileSystem(t *testing.T) {
	fs := &FileSystem{}

	if !fs.DirectoryExists("/tmp") {
		t.Errorf("Expect /tmp directory exists.")
	}

	text := "Hello World!"
	dir := "/tmp/goduckdb"

	fs.CreateDirectory(dir)
	handle := fs.OpenFile(dir+fs.PathSeparator()+"foo.txt", WriteOnly|Create, NoLock)
	buffer := NewFileBuffer(4096)
	for i := range text {
		buffer.Buffer()[i] = text[i]
	}
	buffer.Write(handle, 0)

	buffer.Clear()
	buffer.Read(handle, 0)

	for i := range text {
		if b := buffer.Buffer()[i]; b != text[i] {
			t.Errorf("Expect %d, got %d", text[i], b)
		}
	}

	handle.Close()
}
