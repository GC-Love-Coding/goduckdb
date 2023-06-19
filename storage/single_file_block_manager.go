package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/goduckdb/common"
)

const BlockStart = HeaderSize * 3

// SingleFileBlockManager is a implementation for a BlockManager which manages blocks in a single file.
type SingleFileBlockManager struct {
	activeHeader   uint8             // The active DatabaseHeader, either 0 (h1) or 1 (h2).
	path           string            // The path where the file is stored.
	handle         common.FileHandle // The buffer used to read/write to the headers.
	headerBuffer   *common.FileBuffer
	freeList       []BlockID // The list of free blocks that can be written to currently.
	usedBlocks     []BlockID // The list of blocks that are used by the current block manager.
	metaBlock      BlockID   // The current meta block id.
	maxBlock       BlockID   // The current maximum block id, this id will be given away first after the free_list runs out.
	iterationCount uint64    // The current header iteration count.
}

func NewSingleFileBlockManager(fs *common.FileSystem, path string, readOnly bool, createNew bool) BlockManager {
	var flags common.FileFlags
	var lock common.FileLockType

	if readOnly {
		flags = common.ReadOnly | common.DirectIO
		lock = common.ReadLock
	} else {
		flags = common.WriteOnly | common.DirectIO
		lock = common.WriteLock

		if createNew {
			flags |= common.Create
		}
	}

	// Open the RDBMS handle.
	headerBuffer := common.NewFileBuffer(HeaderSize)
	handle := fs.OpenFile(path, flags, lock)

	if createNew {
		// If we create a new file, we fill the metadata of the file
		// first fill in the new header.
		headerBuffer.Clear()
		binary.LittleEndian.PutUint64(headerBuffer.Buffer(), VersionNo)
		headerBuffer.Write(handle, 0)
		headerBuffer.Clear()

		// Write the database headers.
		// Initialize meta_block and free_list to INVALID_BLOCK because
		// the database file does not contain any actual content yet.
		// header 1.
		databaseHeader := DatabaseHeader{
			Iteration:  0,
			MetaBlock:  InvalidBlock,
			FreeList:   InvalidBlock,
			BlockCount: 0,
		}
		data := DatabaseHeaderToBytes(databaseHeader)
		copy(headerBuffer.Buffer(), data)
		headerBuffer.Write(handle, HeaderSize)

		// header 2.
		databaseHeader.Iteration = 1
		data = DatabaseHeaderToBytes(databaseHeader)
		copy(headerBuffer.Buffer(), data)
		headerBuffer.Write(handle, HeaderSize*2)

		// Ensure that writing to disk is completed before returning
		handle.Sync()

		return &SingleFileBlockManager{
			activeHeader: 1,
			path:         path,
			headerBuffer: headerBuffer,
			handle:       handle,
		}
	} else {
		// Otherwise, we check the metadata of the file.
		headerBuffer.Read(handle, 0)
		mainHeader := BytesToMainHeader(headerBuffer.Buffer())

		if mainHeader.VersionNo != VersionNo {
			panic(fmt.Sprintf("Trying to read a database file with version number %d, but we can only read version %d",
				mainHeader.VersionNo, VersionNo))
		}

		var activeHeader uint8
		// Read the database headers from disk.
		headerBuffer.Read(handle, HeaderSize)
		databaseHeader1 := BytesToDatabaseHeader(headerBuffer.Buffer())
		headerBuffer.Read(handle, HeaderSize*2)
		databaseHeader2 := BytesToDatabaseHeader(headerBuffer.Buffer())

		manager := &SingleFileBlockManager{
			activeHeader: activeHeader,
			path:         path,
			headerBuffer: headerBuffer,
			handle:       handle,
		}

		// Check the header with the highest iteration count.
		if databaseHeader1.Iteration > databaseHeader2.Iteration {
			// h1 is ative header.
			manager.Initialize(databaseHeader1)
		} else {
			// h2 is active header.
			manager.activeHeader = 1
			manager.Initialize(databaseHeader2)
		}

		return manager
	}
}

func (manager *SingleFileBlockManager) Initialize(header DatabaseHeader) {
	if header.FreeList != InvalidBlock {
		reader := NewMetaBlockReader(manager, header.FreeList)
		freeListCount := reader.Read(uint64(0)).(uint64)

		for i := 0; uint64(i) < freeListCount; i++ {
			manager.freeList = append(manager.freeList, reader.Read(BlockID(0)).(BlockID))
		}
	}

	manager.metaBlock = header.MetaBlock
	manager.iterationCount = header.Iteration
	manager.maxBlock = BlockID(header.BlockCount)
}

func (manager *SingleFileBlockManager) CreateBlock() *Block {
	bid := manager.GetFreeBlockID()

	return NewBlock(bid)
}

func (manager *SingleFileBlockManager) GetFreeBlockID() BlockID {
	var blockID BlockID

	if size := len(manager.freeList); size > 0 {
		blockID = manager.freeList[size-1]
		manager.freeList = manager.freeList[:size-1]
	} else {
		blockID = manager.maxBlock
		manager.maxBlock++
	}

	return blockID
}

func (manager *SingleFileBlockManager) GetMetaBlock() BlockID {
	return manager.metaBlock
}

func (blockManager *SingleFileBlockManager) Read(block *Block) {
	// TODO: duplicate block ids
	blockManager.usedBlocks = append(blockManager.usedBlocks, block.ID)
	block.Read(blockManager.handle, uint64(BlockStart+block.ID*BlockSize))
}

func (blockManager *SingleFileBlockManager) Write(block *Block) {
	block.Write(blockManager.handle, uint64(BlockStart+block.ID*BlockSize))
}

// TODO: how it works?
func (manager *SingleFileBlockManager) WriteHeader(header DatabaseHeader) {
	// Set the iteration count.
	header.Iteration = manager.iterationCount
	manager.iterationCount++

	// Now handle the free list.
	if len(manager.usedBlocks) > 0 {
		// There are blocks in the free list.
		// Write them to the file.
		writer := NewMetaBlockWriter(manager)
		header.FreeList = writer.block.ID
		writer.Write(len(manager.usedBlocks))

		for _, blockID := range manager.usedBlocks {
			writer.Write(blockID)
		}
		writer.Flush()
	} else {
		// No block in the free list.
		header.FreeList = InvalidBlock
	}

	// Set the header inside the buffer.
	manager.headerBuffer.Clear()
	data := DatabaseHeaderToBytes(header)
	copy(manager.headerBuffer.Buffer(), data)
	// Now write the header to the file, active_header determines whether we write to h1 or h2.
	// Note that if active_header is h1 we write to h2, and vice versa.
	if manager.activeHeader == 1 {
		manager.headerBuffer.Write(manager.handle, HeaderSize)
	} else {
		manager.headerBuffer.Write(manager.handle, HeaderSize*2)
	}
	// Switch active header to the other header.
	manager.activeHeader = 1 - manager.activeHeader
	// Ensure the header to the other header.
	manager.handle.Sync()

	// The free list is now equal to the blocks that were used by previous iteration.
	manager.freeList = manager.usedBlocks
}
