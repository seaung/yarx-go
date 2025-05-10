package tasks

import (
	"github.com/seaung/yarx-go/internal/yarx/service"
	"github.com/seaung/yarx-go/internal/yarx/store"
)

type TasksController struct {
	s service.IService
}

func NewTasksController(store store.IStore) *TasksController {
	return &TasksController{
		s: service.NewService(store),
	}
}
