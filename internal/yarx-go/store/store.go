package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	S    *datastore
)

// IStore接口，定义了Store需要实现的方法
type IStore interface {
	DB() *gorm.DB
	Users() UserStore
	TaskCrons() TaskCronStore
	TaskRuns() TaskRunStore
	TaskMains() TaskMainStore
}

// IStore的一个实体
type datastore struct {
	db *gorm.DB
}

// 确保datastore实现了IStore接口
var _ IStore = (*datastore)(nil)

// NewStore用于创建一个IStore类型实例
func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db: db}
	})
	return S
}

// 实现DB方法
func (ds *datastore) DB() *gorm.DB {
	return ds.db
}

// 实现Users方法返回UserStore接口的实例
func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}

// 实现TaskCronStore方法返回TaskCronStore实例
func (ds *datastore) TaskCrons() TaskCronStore {
	return newTaskCron(ds.db)
}

func (ds *datastore) TaskRuns() TaskRunStore {
	return newTaskRuns(ds.db)
}

func (ds *datastore) TaskMains() TaskMainStore {
	return newTaskMain(ds.db)
}
