package storage

type BlockManager interface {
	// Creates a new block inside the block manager.
	CreateBlock() *Block
	// Return the next free block id.
	GetFreeBlockID() BlockID
	// Get the first meta block id.
	GetMetaBlock() BlockID
	// Read the content of the block from disk.
	Read(block *Block)
	// Writes the block to disk.
	Write(block *Block)
	// Write the header; should be the final step of a checkpoint.
	WriteHeader(header DatabaseHeader)
}
