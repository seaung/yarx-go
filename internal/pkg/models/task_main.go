package models

import "time"

type TaskMain struct {
	ID              int64     `gorm:"column:id;primary_key"`
	TaskID          string    `gorm:"column:task_id"`
	TaskName        string    `gorm:"column:task_name"`
	KwArgs          string    `gorm:"column:kwargs"`
	State           string    `gorm:"column:state"`
	Result          string    `gorm:"column:result"`
	CreatedDatetime time.Time `gorm:"column:created_datetime"`
	UpdatedDatetime time.Time `gorm:"column:updated_datetime"`
}

func (t *TaskMain) TableName() string {
	return "task_main"
}
