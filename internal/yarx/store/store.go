package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once  sync.Once
	Store *datastore
)

type datastore struct {
	ds *gorm.DB
}

func NewStore(ds *gorm.DB) *datastore {
	once.Do(func() {
		Store = &datastore{ds: ds}
	})

	return Store
}
