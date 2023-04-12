package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type TaskCronStore interface {
	CreateTask(ctx context.Context, task *models.TaskCron) error
	GetTaskByID(ctx context.Context, id int64) error
	GetTaskByTaskID(ctx context.Context, taskID int64) error
	UpdateTask(ctx context.Context, tasks *models.TaskCron) error
	DeleteTask(ctx context.Context, taskID int64) error
}

type taskCron struct {
	db *gorm.DB
}

var _ TaskCronStore = (*taskCron)(nil)

func newTaskCron(db *gorm.DB) *taskCron {
	return &taskCron{db}
}

func (t *taskCron) CreateTask(ctx context.Context, task *models.TaskCron) error {
	return nil
}

func (t *taskCron) GetTaskByID(ctx context.Context, id int64) error {
	return nil
}

func (t *taskCron) GetTaskByTaskID(ctx context.Context, taskID int64) error {
	return nil
}

func (t *taskCron) UpdateTask(ctx context.Context, tasks *models.TaskCron) error {
	return nil
}

func (t *taskCron) DeleteTask(ctx context.Context, taskID int64) error {
	return nil
}
