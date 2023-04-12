package models

import "time"

type TaskCron struct {
	ID              int64     `gorm:"column:id;primary_key"`
	TaskID          string    `gorm:"column:task_id"`
	TaskName        string    `gorm:"column:task_name"`
	KwArgs          string    `gorm:"column:kwargs"`
	CronRule        string    `gorm:"column:cron_rule"`
	Status          string    `gorm:"column:status"`
	Description     string    `gorm:"column:description"`
	RunCount        int       `gorm:"column:run_count"`
	WorkspaceID     int       `gorm:"column:worksapce_id"`
	CreatedDatetime time.Time `gorm:"column:created_datetime"`
	UpdatedDatetime time.Time `gorm:"column:updated_datetime"`
	LastRunDatetime time.Time `gorm:"column:last_run_datetime"`
}

func (t *TaskCron) TableName() string {
	return "task_cron"
}
