package common

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

const FileBufferBlockSize = 4096
const FileBufferHeaderSize = uint64(unsafe.Sizeof(uint64(0)))

// The FileBuffer represents a buffer that can be read or written to a Direct IO FileHandle.
type FileBuffer struct {
	size         uint64 // The size of the portion that users can write to, this is equivalent to internal_size - FILE_BUFFER_HEADER_SIZE
	internalSize uint64 // The aligned size as passed to the constructor. This is the size that is read or written to disk.
	buffer       []byte // The buffer that users can write to
	// internalBuffer *[]byte // The pointer to the internal buffer that will be read or written, including the buffer header
	// mallocedBuffer []byte  // The buffer that was actually malloc'd, i.e. the pointer that must be freed when the FileBuffer is destroyed
}

// TODO: ignore the alignment
func NewFileBuffer(bufSize uint64) *FileBuffer {
	return &FileBuffer{
		size:         bufSize - FileBufferHeaderSize,
		internalSize: bufSize,
		buffer:       make([]byte, bufSize),
	}
	// mallocedBuffer := make([]byte, bufSize+FileBufferBlockSize-1)
	// unalignedPtr := uint64(uintptr(unsafe.Pointer(&mallocedBuffer[0])))
	// alignedPtr := unalignedPtr

	// if remainder := alignedPtr % FileBufferBlockSize; remainder != 0 {
	// 	alignedPtr += FileBufferBlockSize - remainder
	// }

	// internalBuffer := mallocedBuffer[alignedPtr-unalignedPtr:]
	// internalSize := bufSize
	// buffer := internalBuffer[FileBufferHeaderSize:]
	// size := internalSize - FileBufferHeaderSize

	// return &FileBuffer{
	// 	buffer:         buffer,
	// 	size:           size,
	// 	internalBuffer: internalBuffer,
	// 	internalSize:   bufSize,
	// 	mallocedBuffer: mallocedBuffer,
	// }
}

func (fb *FileBuffer) Buffer() []byte {
	return fb.buffer[FileBufferHeaderSize:]
}

func (fb *FileBuffer) Size() uint64 {
	return fb.size
}

func (fb *FileBuffer) Read(handle FileHandle, offset uint64) {
	handle.Read(fb.buffer, offset)
	storedChecksum := binary.LittleEndian.Uint64(fb.buffer[:FileBufferHeaderSize])
	computedChecksum := Checksum(fb.buffer[FileBufferHeaderSize:])

	if computedChecksum != storedChecksum {
		panic(fmt.Sprintf("Corrupt database file: computed checksum %x does not match stored checksum %x in block", computedChecksum, storedChecksum))
	}
}

func (fb *FileBuffer) Write(handle FileHandle, offset uint64) {
	checksum := Checksum(fb.buffer[FileBufferHeaderSize:])
	binary.LittleEndian.PutUint64(fb.buffer, checksum)
	handle.Write(fb.buffer, offset)
}

func bzero(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

func (fb *FileBuffer) Clear() {
	bzero(fb.buffer)
}
