package tasks

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
)

func (t *taskservice) Get(ctx context.Context, taskID string) (*models.Tasks, error) {
	return t.ds.Tasks().GetTaskByID(ctx, taskID)
}
