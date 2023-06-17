package common

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

type FileLockType uint8
type FileFlags uint8

const (
	NoLock FileLockType = 1 << iota
	ReadLock
	WriteLock
)

const (
	ReadOnly FileFlags = 1 << iota
	WriteOnly
	DirectIO
	Create
)

type FileSystem struct{}

func (fs *FileSystem) OpenFile(path string, flags FileFlags, lockType FileLockType) FileHandle {
	var openFlags int

	if flags&ReadOnly != 0 {
		openFlags = os.O_RDONLY
	} else {
		// TODO: seems no O_CLOEXEC
		// since we don't need to fork, just ignore it temporarily
		openFlags = os.O_RDWR

		if flags&Create != 0 {
			openFlags |= os.O_CREATE
		}
	}

	// OSX does not have O_DIRECT, instead we need to use fcntl afterwards to support direct IO
	if flags&DirectIO != 0 {
		openFlags |= os.O_SYNC
	}

	file, err := os.OpenFile(path, openFlags, 0666)

	if err != nil {
		panic("Cannot open file " + path)
	}

	if flags&DirectIO != 0 {
		_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, file.Fd(), syscall.F_NOCACHE, 1)

		if errno != 0 {
			panic("Could not enable direct IO for file " + path)
		}
	}

	// Set lock on file.
	if lockType != NoLock {
		flock := syscall.Flock_t{
			Type:   syscall.F_RDLCK,
			Whence: io.SeekStart,
			Start:  0,
			Len:    0,
		}

		if lockType == WriteLock {
			flock.Type = syscall.F_WRLCK
		}

		err := syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &flock)

		if err != nil {
			file.Close()
			panic("Could not set lock on file " + path)
		}
	}

	return NewFileHandle(fs, file, path)
}

// Read exactly nbytes from the specified offset in the file. Fails if nbytes could not be read.
// This is equivalent to calling SetFilePointer(offset) followed by calling Read().
func (fs *FileSystem) ReadFromOffset(handle FileHandle, buffer *[]byte, nbyte uint64, offset uint64) {
	fs.SetFilePointer(handle, offset)
	fs.Read(handle, buffer, nbyte)
}

// TODO: consider syscall.Read
func (fs *FileSystem) Read(handle FileHandle, buffer *[]byte, nbyte uint64) int64 {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file

	temp := make([]byte, nbyte)
	bytesRead, err := file.Read(temp)

	if err != nil {
		panic("Could not read from file " + unixHandle.path)
	}

	*buffer = append(*buffer, temp...)

	return int64(bytesRead)
}

// Write nbyte from the buffer into the file, moving the file pointer forward by nbyte.
func (fs *FileSystem) Write(handle FileHandle, buffer []byte, nbyte uint64) int64 {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	n, err := file.Write(buffer[:nbyte])

	if err != nil {
		panic("Could not write file " + unixHandle.path)
	}

	return int64(n)
}

func (fs *FileSystem) WriteFromOffset(handle FileHandle, buffer []byte, nbyte uint64, offset uint64) {
	fs.SetFilePointer(handle, offset)
	bytesWritten := fs.Write(handle, buffer, nbyte)

	if uint64(bytesWritten) != nbyte {
		panic("Could not write sufficient bytes from file " + handle.(*UnixFileHandle).path)
	}
}

// Set the file pointer of a file handle to a specified offset.
// Reads and writes will happen from this location
func (fs *FileSystem) SetFilePointer(handle FileHandle, offset uint64) {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	_, err := file.Seek(int64(offset), os.SEEK_SET)

	if err != nil {
		panic(fmt.Sprintf("Could not seek to location %d for file %s", offset, unixHandle.path))
	}
}
