package service

import (
	"github.com/seaung/yarx-go/internal/yarx/service/users"
	"github.com/seaung/yarx-go/internal/yarx/store"
)

type IService interface {
	Users() users.UserService
}

type service struct {
	ds store.IStore
}

var _ IService = (*service)(nil)

func NewService(ds store.IStore) *service {
	return &service{
		ds: ds,
	}
}

func (s *service) Users() users.UserService {
	return users.NewUserService(s.ds)
}
