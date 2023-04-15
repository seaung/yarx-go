package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type TaskMainStore interface {
	CreateMainTask(ctx context.Context, task *models.TaskMain) error
	GetTaskByID(ctx context.Context, id int64) (*models.TaskMain, error)
	GetTaskMainByName(ctx context.Context, name string) (*models.TaskMain, error)
	UpdateTaskMain(ctx context.Context, task *models.TaskMain) error
	DeleteTaskMain(ctx context.Context, taskID int64) error
}

type taskMain struct {
	db *gorm.DB
}

var _ TaskMainStore = (*taskMain)(nil)

func newTaskMain(db *gorm.DB) *taskMain {
	return &taskMain{db: db}
}

func (t *taskMain) CreateMainTask(ctx context.Context, task *models.TaskMain) error {
	return nil
}

func (t *taskMain) GetTaskByID(ctx context.Context, id int64) (*models.TaskMain, error) {
	return nil, nil
}

func (t *taskMain) GetTaskMainByName(ctx context.Context, name string) (*models.TaskMain, error) {
	return nil, nil
}

func (t *taskMain) UpdateTaskMain(ctx context.Context, task *models.TaskMain) error {
	return nil
}

func (t *taskMain) DeleteTaskMain(ctx context.Context, taskID int64) error {
	return nil
}
