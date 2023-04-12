package biz

import (
	"github.com/seaung/yarx-go/internal/yarx-go/biz/user"
	"github.com/seaung/yarx-go/internal/yarx-go/store"
)

type IBiz interface {
	Users() user.UserBiz
}

var _ IBiz = (*biz)(nil)

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
