package catalog

type CatalogType uint8

const (
	Invalid CatalogType = iota
	Table
	Schema
	TableFunction
	ScalarFunction
	View
	Index
	UpdatedEntry      = 10
	DeletedEntry      = 11
	PreparedStatement = 12
	Sequence          = 13
)
