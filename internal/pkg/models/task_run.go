package models

import "time"

type TaskRun struct {
	ID              int64     `gorm:"column:id;primary_key"`
	TaskID          string    `gorm:"column:task_id"`
	TaskName        string    `gorm:"column:task_name"`
	Kwargs          string    `gorm:"column:kwargs"`
	CreatedDatetime time.Time `gorm:"column:created_datetime"`
	UpdatedDatetime time.Time `gorm:"column:updated_datetime"`
}

func (t *TaskRun) TableName() string {
	return "task_run"
}
