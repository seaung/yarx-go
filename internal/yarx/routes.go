package yarx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func installRoutes(g *gin.Engine) error {
	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "请求的接口不存在"})
	})

	v1 := g.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/login", func(ctx *gin.Context) {})
		}

		tasks := v1.Group("/tasks")
		{
			tasks.POST("/created", func(ctx *gin.Context) {})
		}

	}
	return nil
}
