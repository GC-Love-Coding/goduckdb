package storage

import (
	"encoding/binary"
	"unsafe"
)

// Size of a memory slot managed by the StorageManager. This is the quantum of allocation for Blocks on DuckDB. We
// default to 256KB. (1 << 18)
const (
	BlockSize    = 262144
	HeaderSize   = 4096
	InvalidBlock = -1
	VersionNo    = 1
)

// The MainHeader is the first header in the storage file.
// The MainHeader is typically written only once for a database file.
type MainHeader struct {
	VersionNo uint64 // The version of the database.
	Flags     [4]uint64
}

// The DatabaseHeader contains information about the current state of the database. Every storage file has two
// DatabaseHeaders. On startup, the DatabaseHeader with the highest iteration count is used as the active header. When
// a checkpoint is performed, the active DatabaseHeader is switched by increasing the iteration count of the
// DatabaseHeader.
type DatabaseHeader struct {
	Iteration  uint64  // The iteration count, increases by 1 every time the storage is checkpointed.
	MetaBlock  BlockID // A pointer to the initial meta block.
	FreeList   BlockID // A pointer to the block containing the free list.
	BlockCount uint64  // The number of blocks that is in the file as of this database header. If the file is larger than BLOCK_SIZE * block_count any blocks appearing AFTER block_count are implicitly part of the free_list.
}

type ByteSlice []byte

func (bs *ByteSlice) Write(p []byte) (n int, err error) {
	*bs = append(*bs, p...)

	return len(p), nil
}

func DatabaseHeaderToBytes(header DatabaseHeader) []byte {
	buffer := new(ByteSlice)
	binary.Write(buffer, binary.LittleEndian, header)

	return []byte(*buffer)
}

func BytesToDatabaseHeader(buffer []byte) DatabaseHeader {
	var header DatabaseHeader
	header.Iteration = binary.LittleEndian.Uint64(buffer)
	buffer = buffer[unsafe.Sizeof(header.Iteration):]
	header.MetaBlock = BlockID(binary.LittleEndian.Uint32(buffer))
	buffer = buffer[unsafe.Sizeof(header.MetaBlock):]
	header.FreeList = BlockID(binary.LittleEndian.Uint32(buffer))
	buffer = buffer[unsafe.Sizeof(header.FreeList):]
	header.BlockCount = binary.LittleEndian.Uint64(buffer)

	return header
}

func BytesToMainHeader(buffer []byte) MainHeader {
	var header MainHeader
	header.VersionNo = binary.LittleEndian.Uint64(buffer)
	buffer = buffer[unsafe.Sizeof(header.VersionNo):]

	for i := range header.Flags {
		header.Flags[i] = binary.LittleEndian.Uint64(buffer)
		buffer = buffer[unsafe.Sizeof(header.Flags[i]):]
	}

	return header
}
