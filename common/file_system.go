package common

type FileLockType uint8
type FileFlags uint8

const (
	NoLock FileLockType = iota
	ReadLock
	WriteLock
)

const (
	Read FileFlags = iota
	Write
	DirectIO
	Create
)

type FileSystem struct{}

func (fs FileSystem) OpenFile(path string, flags FileFlags, lock FileLockType) FileHandle {
	return nil
}
