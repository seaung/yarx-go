package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *TasksController) GetTasks(ctx *gin.Context) {
	task_id := ctx.Params.ByName("task_id")
	if task_id == "" {
		ctx.JSON(400, gin.H{"error": "task_id is required"})
		return
	}

	data, err := ctrl.s.Tasks().Get(ctx, task_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "server internal error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "ok", "data": data})
}
