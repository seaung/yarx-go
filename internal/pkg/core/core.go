package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/internal/pkg/errno"
)

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		hcode, code, message := errno.Decode(err)
		c.JSON(hcode, ResponseError{
			Code:    code,
			Message: message,
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
