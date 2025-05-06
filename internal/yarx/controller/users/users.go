package users

import (
	"github.com/seaung/yarx-go/internal/yarx/service"
	"github.com/seaung/yarx-go/internal/yarx/store"
)

type UserController struct {
	s service.IService
}

func NewUserController(s store.IStore) *UserController {
	return &UserController{
		s: service.NewService(s),
	}
}
