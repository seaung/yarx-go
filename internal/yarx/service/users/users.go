package users

import (
	"context"

	"github.com/seaung/yarx-go/internal/yarx/store"
	"github.com/seaung/yarx-go/pkg/api"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*api.LoginResponse, error)
	IsExist(ctx context.Context, username, email string) (bool, error)
	CreateUser(ctx context.Context, form *api.CreateUserRequestForm) error
}

type userService struct {
	ds store.IStore
}

func NewUserService(ds store.IStore) *userService {
	return &userService{
		ds: ds,
	}
}

var _ UserService = (*userService)(nil)
