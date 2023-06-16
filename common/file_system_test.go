package common

import "testing"

func TestFileSystem(t *testing.T) {
	fs := FileSystem{}
	fs.OpenFile("a.txt", Write, NoLock)
}
