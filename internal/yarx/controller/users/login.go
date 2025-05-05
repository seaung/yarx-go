package users

import (
	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/pkg/api"
)

func (ctrl *UserController) Login(ctx *gin.Context) {
	var req api.LoginRequestForm

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
	}

	resp, err := ctrl.s.Users().Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "登录失败",
		})
	}

	ctx.JSON(200, resp)
	return
}
