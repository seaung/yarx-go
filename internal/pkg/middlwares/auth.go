package middlwares

import (
	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/pkg/constants"
	"github.com/seaung/yarx-go/internal/pkg/core"
	"github.com/seaung/yarx-go/internal/pkg/errno"
	"github.com/seaung/yarx-go/pkg/token"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.TokenInvalidError, nil)
			c.Abort()
			return
		}
		c.Set(constants.XUsernameKey, username)
		c.Next()
	}
}
