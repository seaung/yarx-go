package tasks

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"github.com/seaung/yarx-go/internal/yarx/store"
)

type TasksService interface {
	Create(ctx context.Context) error
	Get(ctx context.Context, taskID string) (*models.Tasks, error)
}

type taskservice struct {
	ds store.IStore
}

func NewTaskService(ds store.IStore) *taskservice {
	return &taskservice{ds}
}

var _ TasksService = (*taskservice)(nil)
