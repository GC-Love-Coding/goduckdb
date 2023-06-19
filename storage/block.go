package storage

import "github.com/goduckdb/common"

type BlockID int64

type Block struct {
	*common.FileBuffer
	ID BlockID
}

// TODO: return *Block or Block
func NewBlock(id BlockID) *Block {
	return &Block{
		FileBuffer: common.NewFileBuffer(BlockSize),
		ID:         id,
	}
}
