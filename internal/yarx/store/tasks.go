package store

import (
	"context"

	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type tasks struct {
	db *gorm.DB
}

type TasksStore interface {
	GetTaskByID(ctx context.Context, taskID string) (*models.Tasks, error)
}

func newTasks(ds *gorm.DB) *tasks {
	return &tasks{db: ds}
}

func (t *tasks) GetTaskByID(ctx context.Context, taskID string) (*models.Tasks, error) {
	var data models.Tasks
	err := t.db.Table("tasks").Select("id,tasks_id,task_name,task_type").Where("task_id = ?", taskID).First(&data).Error
	if err != nil {
		return &models.Tasks{}, err
	}

	return &data, nil
}
