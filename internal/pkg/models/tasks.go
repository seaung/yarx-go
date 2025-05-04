package models

import "time"

type Tasks struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	TaskID     string    `gorm:"column:task_id"`
	TaskName   string    `gorm:"column:task_name"`
	TaskType   string    `gorm:"column:task_type"`
	TaskStatus string    `gorm:"column:task_status"`
	Targets    string    `gorm:"column:targets"`
	NextTime   time.Time `gorm:"column:next_time"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (t *Tasks) TableName() string {
	return "tasks"
}

type TastkResult struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	TaskID     string    `gorm:"column:task_id"`
	Target     string    `gorm:"column:target"`
	Header     string    `gorm:"column:header"`
	Request    string    `gorm:"column:request"`
	Response   string    `gorm:"column:response"`
	StatusCode int       `gorm:"column:status_code"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (t *TastkResult) TableName() string {
	return "task_results"
}
