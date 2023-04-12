package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"
	"github.com/seaung/yarx-go/internal/pkg/errno"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"github.com/seaung/yarx-go/internal/yarx-go/store"
	v1 "github.com/seaung/yarx-go/pkg/api/yarx-go/v1"
	"github.com/seaung/yarx-go/pkg/auth"
	"github.com/seaung/yarx-go/pkg/token"
)

type UserBiz interface {
	Login(ctx context.Context, r *v1.LoginRequestForm) (*v1.LoginResponse, error)
	Create(ctx context.Context, r *v1.CreateUserRequestForm) error
}

// userBiz接口的实现
type userBiz struct {
	db store.IStore
}

// 确保userBiz实现了UserBiz接口
var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{db: ds}
}

// 实现UserBiz接口中的Login方法
func (u *userBiz) Login(ctx context.Context, r *v1.LoginRequestForm) (*v1.LoginResponse, error) {
	// 1. 获取登录用户的所有信息
	user, err := u.db.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.UserNotFoundError
	}

	// 2. 判断前端发送过来的密码和数据库中的用户密码是否一致
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.PasswordIncorrectError
	}

	// 3. 签发token
	tk, err := token.SignToken(r.Username)
	if err != nil {
		return nil, errno.SignTokenError
	}

	// 4. 成功返回token
	return &v1.LoginResponse{Token: tk}, nil
}

func (u *userBiz) Create(ctx context.Context, r *v1.CreateUserRequestForm) error {
	var user models.UserModel

	_ = copier.Copy(&user, r)

	if err := u.db.Users().Create(ctx, &user); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.UserAlreadyExistError
		}

		return err
	}

	return nil
}
