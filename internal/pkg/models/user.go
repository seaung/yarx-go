package models

import (
	"time"

	"github.com/seaung/yarx-go/pkg/auth"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       int64     `gorm:"column:id;primary_key"`
	Username string    `gorm:"column:username;not null"`
	Password string    `gorm:"column:password;not null"`
	NickName string    `gorm:"column:nickname"`
	Email    string    `gorm:"column:email"`
	CreateAt time.Time `gorm:"column;createAt"`
	UpdateAt time.Time `gorm:"column;updateAt"`
}

func (u *UserModel) TableName() string {
	return "users"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}

	return nil
}
