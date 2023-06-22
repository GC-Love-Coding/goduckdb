package catalog

// TODO: 弄明白child和parent的含义
type CatalogEntry struct {
	ctype     CatalogType   // The type of this catalog entry.
	catalog   *Catalog      // Reference to the catalog this entry belongs to.
	set       *CatalogSet   // Reference to the catalog set this entry is stored in.
	name      string        // The name of the entry.
	deleted   bool          // Whether or not the object is deleted.
	timestamp uint64        // Timestamp at which the catalog entry was created.
	child     *CatalogEntry // Child entry.
	parent    *CatalogEntry // Parent entry (the node that owns this node).
}

func NewCatalogEntry(ctype CatalogType, catalog *Catalog, name string) *CatalogEntry {
	return &CatalogEntry{
		ctype:ctype,
		catalog: catalog
		name: name,
	}
}