package users

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jinzhu/copier"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"github.com/seaung/yarx-go/pkg/api"
)

func (u *userService) IsExist(ctx context.Context, username, email string) (bool, error) {
	return u.ds.Users().IsExist(ctx, username, email)
}

func (u *userService) CreateUser(ctx context.Context, form *api.CreateUserRequestForm) error {
	var account models.Account
	_ = copier.Copy(&account, form)

	if err := u.ds.Users().CreateUser(ctx, &account); err != nil {
		if ok, _ := regexp.MatchString("Duplicate entry", err.Error()); ok {
			return fmt.Errorf("用户名或邮箱已存在")
		}
		return err
	}
	return nil
}
