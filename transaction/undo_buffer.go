package transaction

type UndoFlags uint8

const (
	EmptyEntry UndoFlags = iota
	CatalogEntry
	InsertTuple
	DeleteTuple
	UpdateTuple
	Query
)

type UndoEntry struct {
	utype  UndoFlags
	length uint64
	data   []byte
}

// The undo buffer of a transaction is used to hold previous versions of tuples
// that might be required in the future (because of rollbacks or previous
// transactions accessing them).
type UndoBuffer struct {
	entries []UndoEntry // List of UndoEntries, FIXME: this can be more efficient.
}
