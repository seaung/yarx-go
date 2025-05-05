package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once  sync.Once
	Store *datastore
)

// interface for store
type IStore interface {
	DB() *gorm.DB
	Users() UserStore
}

type datastore struct {
	ds *gorm.DB
}

func NewStore(ds *gorm.DB) *datastore {
	once.Do(func() {
		Store = &datastore{ds: ds}
	})

	return Store
}

var _ IStore = (*datastore)(nil)

func (s *datastore) DB() *gorm.DB {
	return s.ds
}

func (s *datastore) Users() UserStore {
	return newUsers(s.ds)
}
