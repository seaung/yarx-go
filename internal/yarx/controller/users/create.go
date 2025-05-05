package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/pkg/api"
)

func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	var req api.CreateUserRequestForm

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "bad request"})
	}

	if ok, err := ctrl.s.Users().IsExist(ctx, req.Nickname, req.Email); err != nil && !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "internal server error"})
	}

	if err := ctrl.s.Users().CreateUser(ctx, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "internal server error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok"})
}
