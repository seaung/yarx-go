package yarx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/yarx/controller/users"
	"github.com/seaung/yarx-go/internal/yarx/store"
)

func installRoutes(g *gin.Engine) error {
	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "请求的接口不存在"})
	})

	uc := users.NewUserController(store.Store)

	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/created", uc.CreateUser)
		}

		tasks := v1.Group("/tasks")
		{
			tasks.POST("/created", func(ctx *gin.Context) {})
		}

	}
	return nil
}
