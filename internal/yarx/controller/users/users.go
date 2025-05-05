package users

import "github.com/seaung/yarx-go/internal/yarx/service"

type UserController struct {
	s service.IService
}

func NewUserController() *UserController {
	return &UserController{}
}
