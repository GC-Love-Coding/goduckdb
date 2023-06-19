package storage

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// This struct is responsible for writing metadata to disk.
type MetaBlockWriter struct {
	manager BlockManager
	block   *Block
	offset  uint64
}

func NewMetaBlockWriter(manager BlockManager) *MetaBlockWriter {
	return &MetaBlockWriter{manager: manager, block: manager.CreateBlock(), offset: uint64(unsafe.Sizeof(BlockID(0)))}
}

func (writer *MetaBlockWriter) Flush() {
	if writer.offset > uint64(unsafe.Sizeof(BlockID(0))) {
		writer.manager.Write(writer.block)
		writer.offset = uint64(unsafe.Sizeof(BlockID(0)))
	}
}

// Note: offset is updated in `Flush` method.
func (writer *MetaBlockWriter) WriteData(buffer []byte) {
	for writer.offset+uint64(len(buffer)) > writer.block.Size() {
		// We need to make a new block.
		// First copy what we can.
		if copyAmount := writer.block.Size() - writer.offset; copyAmount > 0 {
			copy(writer.block.Buffer()[writer.offset:], buffer[:copyAmount])
			buffer = buffer[copyAmount:]
			writer.offset += copyAmount
		}

		// Now we need to get a new block id.
		newBlockID := writer.manager.GetFreeBlockID()
		// Write the block id of the new block to the start of current block.
		binary.LittleEndian.PutUint64(writer.block.Buffer(), uint64(newBlockID))
		// First flush the old block.
		writer.Flush()
		// Now update the block id of the block.
		writer.block.ID = newBlockID
	}

	copy(writer.block.Buffer(), buffer)
	writer.offset += uint64(len(buffer))
}

func (writer *MetaBlockWriter) Write(v interface{}) {
	switch v.(type) {
	case uint64:
		writer.writeUint64(v.(uint64))
	case int64:
		writer.writeInt64(v.(int64))
	case uint32:
		writer.writeUint32(v.(uint32))
	case int32:
		writer.writeInt32(v.(int32))
	case uint16:
		writer.writeUint16(v.(uint16))
	case int16:
		writer.writeInt16(v.(int16))
	case uint8:
		writer.writeUint8(v.(uint8))
	case int8:
		writer.writeInt8(v.(int8))
	default:
		panic(fmt.Sprintf("Unknown type: %T", v))
	}
}

func (writer *MetaBlockWriter) writeUint64(v uint64) {
	buffer := make([]byte, unsafe.Sizeof(v))
	binary.LittleEndian.PutUint64(buffer, v)

	writer.WriteData(buffer)
}

func (writer *MetaBlockWriter) writeInt64(v int64) {
	writer.writeUint64(uint64(v))
}

func (writer *MetaBlockWriter) writeUint32(v uint32) {
	buffer := make([]byte, unsafe.Sizeof(v))
	binary.LittleEndian.PutUint32(buffer, v)

	writer.WriteData(buffer)
}

func (writer *MetaBlockWriter) writeInt32(v int32) {
	writer.writeUint32(uint32(v))
}

func (writer *MetaBlockWriter) writeUint16(v uint16) {
	buffer := make([]byte, unsafe.Sizeof(v))
	binary.LittleEndian.PutUint16(buffer, v)

	writer.WriteData(buffer)
}

func (writer *MetaBlockWriter) writeInt16(v int16) {
	writer.writeUint16(uint16(v))
}

func (writer *MetaBlockWriter) writeUint8(v uint8) {
	writer.WriteData([]byte{v})
}

func (writer *MetaBlockWriter) writeInt8(v int8) {
	writer.writeUint8(uint8(v))
}
