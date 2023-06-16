package common

type FileHandle interface {
	Read(buffer []byte, nbytes uint64, offset uint64)
	Write(buffer []byte, nbytes uint64, offset uint64)
	Sync()
	Close()
}

type UnixFileHandle struct{}

func NewFileHandle() FileHandle {
	return nil
}
