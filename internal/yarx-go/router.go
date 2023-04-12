package yarxgo

import (
	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/pkg/core"
	"github.com/seaung/yarx-go/internal/pkg/errno"
	"github.com/seaung/yarx-go/internal/pkg/middlwares"
	"github.com/seaung/yarx-go/internal/yarx-go/controllers/v1/user"
	"github.com/seaung/yarx-go/internal/yarx-go/store"
)

func initRouters(g *gin.Engine) error {
	// 注册404 handler
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.PageNotFoundError, nil)
	})

	// 实例化user controller
	uc := user.New(store.S)

	g.POST("/login", uc.Login)

	// 创建路由分组
	v1 := g.Group("/v1", middlwares.AuthToken())
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("/create", uc.Login)
		}
	}

	return nil
}
