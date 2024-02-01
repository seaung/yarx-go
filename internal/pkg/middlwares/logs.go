package middlwares

import (
	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/pkg/log"
)

func LogginMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method
        rePath := c.Request.URL.Path

        remoteIP := c.Request.RemoteAddr

        log.C(c).Info("当前请求", "请求方法", method, "请求路径", rePath, "请求IP", remoteIP)
    }
}
