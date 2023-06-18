package common

import "os"

type FileHandle interface {
	Read(buffer []byte, offset uint64)
	Write(buffer []byte, offset uint64)
	Sync()
	Close()
}

type UnixFileHandle struct {
	*FileSystem
	file *os.File
	path string
}

func NewFileHandle(fs *FileSystem, file *os.File, path string) *UnixFileHandle {
	return &UnixFileHandle{
		FileSystem: fs,
		file:       file,
		path:       path,
	}
}

func (handle *UnixFileHandle) Read(buffer []byte, offset uint64) {
	handle.ReadFromOffset(handle, buffer, offset)
}

func (handle *UnixFileHandle) Write(buffer []byte, offset uint64) {
	handle.WriteFromOffset(handle, buffer, offset)
}

func (handle *UnixFileHandle) Sync() {
	handle.FileSync(handle)
}

func (handle *UnixFileHandle) Close() {
	handle.file.Close()
}
