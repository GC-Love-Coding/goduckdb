package common

import (
	"unsafe"
)

const FileBufferBlockSize = 4096
const FileBufferHeaderSize = uint64(unsafe.Sizeof(uint64(0)))

// The FileBuffer represents a buffer that can be read or written to a Direct IO FileHandle.
type FileBuffer struct {
	buffer         []byte // The buffer that users can write to
	size           uint64 // The size of the portion that users can write to, this is equivalent to internal_size - FILE_BUFFER_HEADER_SIZE
	internalBuffer []byte // The pointer to the internal buffer that will be read or written, including the buffer header
	internalSize   uint64 // The aligned size as passed to the constructor. This is the size that is read or written to disk.
	mallocedBuffer []byte // The buffer that was actually malloc'd, i.e. the pointer that must be freed when the FileBuffer is destroyed
}

// TODO: will it be faster?
func NewFileBuffer(bufSize uint64) *FileBuffer {
	mallocedBuffer := make([]byte, bufSize+FileBufferBlockSize-1)
	unalignedPtr := uint64(uintptr(unsafe.Pointer(&mallocedBuffer[0])))
	alignedPtr := unalignedPtr

	if remainder := alignedPtr % FileBufferBlockSize; remainder != 0 {
		alignedPtr += FileBufferBlockSize - remainder
	}

	internalBuffer := mallocedBuffer[alignedPtr-unalignedPtr:]
	internalSize := bufSize
	buffer := internalBuffer[FileBufferHeaderSize:]
	size := internalSize - FileBufferHeaderSize

	return &FileBuffer{
		buffer:         buffer,
		size:           size,
		internalBuffer: internalBuffer,
		internalSize:   bufSize,
		mallocedBuffer: mallocedBuffer,
	}
}
