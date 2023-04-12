package user

import (
	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/pkg/core"
	"github.com/seaung/yarx-go/internal/pkg/errno"
	"github.com/seaung/yarx-go/internal/pkg/log"
	"github.com/seaung/yarx-go/internal/yarx-go/biz"
	"github.com/seaung/yarx-go/internal/yarx-go/store"
	v1 "github.com/seaung/yarx-go/pkg/api/yarx-go/v1"
)

type UserController struct {
	b biz.IBiz
}

func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}

func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Info("Login function called")

	var r v1.LoginRequestForm
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.BindParameterError, nil)
		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
