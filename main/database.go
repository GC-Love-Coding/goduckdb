package main

import (
	"github.com/goduckdb/common"
	"github.com/goduckdb/storage"
)

type AccessMode int

const (
	Undefined AccessMode = iota
	ReadOnly
	ReadWrite
)

type DBConfig struct {
	accessMode AccessMode
	fileSystem *common.FileSystem
}

// The database object. This object holds the catalog and all the
// database-specific meta information.
type DuckDB struct {
	fileSystem *common.FileSystem
	storage    *storage.StorageManager
}
