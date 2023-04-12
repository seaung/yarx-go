package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type UserStore interface {
	Create(ctx context.Context, user *models.UserModel) error
	Get(ctx context.Context, username string) (*models.UserModel, error)
	Update(ctx context.Context, user *models.UserModel) error
	Delete(ctx context.Context, username string) error
}

type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

func (u *users) Create(ctx context.Context, user *models.UserModel) error {
	return u.db.Create(&user).Error
}

func (u *users) Get(ctx context.Context, username string) (*models.UserModel, error) {
	var user models.UserModel
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *users) Update(ctx context.Context, user *models.UserModel) error {
	return u.db.Save(user).Error
}

func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&models.UserModel{}).Error
	if err != nil {
		return err
	}

	return nil
}
