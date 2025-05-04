package models

import (
	"time"

	"github.com/seaung/yarx-go/pkg/auth"
	"gorm.io/gorm"
)

type Account struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Nickname  string    `gorm:"column:nickname;not null"`
	Password  string    `gorm:"column:password;not null"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (a *Account) TableName() string {
	return "account"
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.Password, err = auth.Encrypt(a.Password)
	if err != nil {
		return err
	}

	return nil
}
