package catalog

import "sync"

type CatalogSet struct {
	catalog     *Catalog
	catalogLock sync.Mutex               // The catalog lock is used to make changes to the data.
	data        map[string]*CatalogEntry // The set of entries present in the CatalogSet.
}

func NewCatalogSet(catalog *Catalog) *CatalogSet {
	return &CatalogSet{
		catalog:     catalog,
		catalogLock: sync.Mutex{},
		data:        make(map[string]*CatalogEntry),
	}
}
