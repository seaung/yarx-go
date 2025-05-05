package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"github.com/seaung/yarx-go/pkg/api"
	"gorm.io/gorm"
)

type users struct {
	ds *gorm.DB
}

type UserStore interface {
	Login(ctx context.Context, form *api.LoginRequestForm) error
	IsExist(ctx context.Context, username, email string) (bool, error)
	CreateUser(ctx context.Context, form *models.Account) error
	GetUser(ctx context.Context, email string) (*models.Account, error)
}

func newUsers(ds *gorm.DB) *users {
	return &users{
		ds: ds,
	}
}

func (u *users) Login(ctx context.Context, form *api.LoginRequestForm) error {
	return nil
}

func (u *users) IsExist(ctx context.Context, username, email string) (bool, error) {
	var exists bool
	err := u.ds.Table("account").
		Select("1").
		Where("nickname = ? or email = ?", username, email).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *users) CreateUser(ctx context.Context, form *models.Account) error {
	return u.ds.Table("account").Create(form).Error
}

func (u *users) GetUser(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	err := u.ds.Table("account").
		Where("email = ?", email).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
