package storage

// The DataBlock is the physical unit to store data it has a physical block which is stored in a file with multiple
// blocks.
// TODO: no definition
type DataBlock struct{}

// Stores the header of each data block.
type BlockHeader struct {
	BlockID        BlockID
	AmountOfTuples uint64
}

// The Block stored in a data block.
// type Block struct {
// 	Header    *BlockHeader
// 	BlockSize uint64     // Block size in Bytes.
// 	Offsets   [10]uint64 // The offset of each column data TODO: define number of columns based on chunck info.
// 	data      []byte
// }
