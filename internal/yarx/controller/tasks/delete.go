package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *TasksController) Delete(ctx *gin.Context) {
	taskId := ctx.Param("task_id")
	if taskId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task_id is empty"})
	}
}
