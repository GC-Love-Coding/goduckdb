package storage

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type MetaBlockReader struct {
	manager   BlockManager
	block     *Block
	offset    uint64
	nextBlock BlockID
}

func NewMetaBlockReader(manager BlockManager, blockID BlockID) *MetaBlockReader {
	reader := &MetaBlockReader{
		manager:   manager,
		block:     NewBlock(-1),
		offset:    0,
		nextBlock: -1,
	}

	reader.readNewBlock(blockID)

	return reader
}

// Read content of size read_size into the buffer.
func (reader *MetaBlockReader) ReadData(outBuffer []byte) {
	inBuffer := reader.block.Buffer()

	for reader.offset+uint64(len(outBuffer)) > reader.block.Size() {
		// Cannot read entire entry from block.
		// First read what we can from this block.
		if toRead := reader.block.Size() - reader.offset; toRead > 0 {
			copy(outBuffer, inBuffer[reader.offset:reader.offset+toRead])
			outBuffer = outBuffer[toRead:]
		}

		// Then move to the next block.
		reader.readNewBlock(reader.nextBlock)
	}

	// We have enough left in this block to read from the buffer.
	copy(outBuffer, inBuffer[reader.offset:])
	reader.offset += uint64(len(outBuffer))
}

func (reader *MetaBlockReader) readNewBlock(blockID BlockID) {
	reader.block.ID = blockID
	reader.manager.Read(reader.block)
	reader.nextBlock = BlockID(binary.LittleEndian.Uint64(reader.block.Buffer()))
	reader.offset = uint64(unsafe.Sizeof(BlockID(0)))
}

func (reader *MetaBlockReader) Read(v interface{}) interface{} {
	buffer := make([]byte, unsafe.Sizeof(v))
	reader.ReadData(buffer)

	switch v.(type) {
	case uint8:
		return buffer[0]
	case uint16:
		return binary.LittleEndian.Uint16(buffer)
	case uint32:
		return binary.LittleEndian.Uint32(buffer)
	case uint64:
		return binary.LittleEndian.Uint64(buffer)
	case int8:
		return int8(buffer[0])
	case int16:
		return int16(binary.LittleEndian.Uint16(buffer))
	case int32:
		return int32(binary.LittleEndian.Uint32(buffer))
	case int64:
		return int64(binary.LittleEndian.Uint64(buffer))

	default:
		panic(fmt.Sprintf("Unknown type: %T", v))
	}
}
