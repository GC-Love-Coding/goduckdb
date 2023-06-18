package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
func (fs *FileSystem) ReadFromOffset(handle FileHandle, buffer []byte, offset uint64) {
	fs.SetFilePointer(handle, offset)
	fs.Read(handle, buffer)
}

// TODO: consider syscall.Read
func (fs *FileSystem) Read(handle FileHandle, buffer []byte) int64 {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	bytesRead, err := file.Read(buffer)

	if err != nil {
		panic("Could not read from file " + unixHandle.path)
	}

	return int64(bytesRead)
}

// Write nbyte from the buffer into the file, moving the file pointer forward by nbyte.
func (fs *FileSystem) Write(handle FileHandle, buffer []byte) int64 {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	n, err := file.Write(buffer)

	if err != nil {
		panic("Could not write file " + unixHandle.path)
	}

	return int64(n)
}

func (fs *FileSystem) WriteFromOffset(handle FileHandle, buffer []byte, offset uint64) {
	fs.SetFilePointer(handle, offset)
	bytesWritten := fs.Write(handle, buffer)

	if bytesWritten != int64(len(buffer)) {
		panic("Could not write sufficient bytes from file " + handle.(*UnixFileHandle).path)
	}
}

// Returns the file size of a file handle, returns -1 on error
func (fs *FileSystem) GetFileSize(handle FileHandle) int64 {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	fileInfo, err := file.Stat()

	if err != nil {
		return -1
	}

	return fileInfo.Size()
}

// Check if a directory exists.
func (fs *FileSystem) DirectoryExists(directory string) bool {
	fileInfo, err := os.Stat(directory)

	if err != nil {
		// if os.IsNotExist(err) {
		// 	return false
		// }
		return false
	}

	return fileInfo.IsDir()
}

// Create a directory if it does not exist.
func (fs *FileSystem) CreateDirectory(directory string) {
	if !fs.DirectoryExists(directory) {
		err := os.Mkdir(directory, 0755)

		if err != nil {
			panic("Failed create directory!")
		}
	}
}

func (fs *FileSystem) RemoveDirectory(directory string) {
	err := os.RemoveAll(directory)

	if err != nil {
		fmt.Println(err)
	}
}

// List files in a directory, invoking the callback method for each one
// TODO: do we need to callback directory?
func (fs *FileSystem) ListFiles(directory string, callback func(string)) bool {
	if !fs.DirectoryExists(directory) {
		return false
	}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			callback(info.Name())
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return true
}

// Move a file from source path to the target, StorageManager relies on this being an atomic action for ACID
// properties
func (fs *FileSystem) MoveFile(source string, target string) {
	// TODO: FIXME: rename does not guarantee atomicity or overwriting target file if it exists
	err := os.Rename(source, target)

	if err != nil {
		panic(err)
	}
}

func (fs *FileSystem) FileExists(filename string) bool {
	_, err := os.Stat(filename)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	panic(err)
}

func (fs *FileSystem) RemoveFile(filename string) {
	err := os.Remove(filename)

	if err != nil {
		panic(err)
	}
}

// Path separator for the current file system.
func (fs *FileSystem) PathSeparator() string {
	return string(filepath.Separator)
}

// Join two paths together.
func (fs *FileSystem) JoinPath(a string, b string) string {
	return filepath.Join(a, b)
}

// Sync a file handle to disk.
func (fs *FileSystem) FileSync(handle FileHandle) {
	unixHandle := handle.(*UnixFileHandle)
	file := unixHandle.file
	file.Sync()
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
