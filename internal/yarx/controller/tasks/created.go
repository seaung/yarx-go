package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/pkg/api"
	"github.com/seaung/yarx-go/pkg/utils"
)

func (ctrl *TasksController) Create(ctx *gin.Context) {
	var req api.CreateTaskRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if req.TaskType == utils.ExecuteImmediately {
		go func() {}()
	}
}
