package users

import (
	"context"
	"fmt"

	"github.com/seaung/yarx-go/pkg/api"
	"github.com/seaung/yarx-go/pkg/auth"
)

func (u *userService) Login(ctx context.Context, email, password string) (*api.LoginResponse, error) {
	account, err := u.ds.Users().GetUser(ctx, email)

	if err != nil {
		return nil, err
	}

	if err = auth.CompareHash(account.Password, password); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(account.Email)

	if err != nil {
		return nil, err
	}

	return &api.LoginResponse{
		UUID:        fmt.Sprintf("uuid-%v", account.ID),
		Email:       account.Email,
		Nickname:    account.Nickname,
		AccessToken: token,
	}, nil
}
