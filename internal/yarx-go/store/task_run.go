package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type TaskRunStore interface {
	Create(ctx context.Context, task *models.TaskRun) error
	Get(ctx context.Context, taskID string) (*models.TaskRun, error)
	Update(ctx context.Context, task *models.TaskRun) error
	Delete(ctx context.Context, taskID string) error
}

type taskRuns struct {
	ds *gorm.DB
}

var _ TaskRunStore = (*taskRuns)(nil)

func newTaskRuns(db *gorm.DB) *taskRuns {
	return &taskRuns{ds: db}
}

func (t *taskRuns) Create(ctx context.Context, task *models.TaskRun) error {
	return nil
}

func (t *taskRuns) Get(ctx context.Context, taskID string) (*models.TaskRun, error) {
	return nil, nil
}

func (t *taskRuns) Update(ctx context.Context, task *models.TaskRun) error {
	return nil
}
func (t *taskRuns) Delete(ctx context.Context, taskID string) error {
	return nil
}
