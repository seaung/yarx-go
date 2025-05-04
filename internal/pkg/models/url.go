package models

import (
	"time"

	"gorm.io/gorm"
)

// CrawlerURL 存储爬虫爬取到的URL信息
type CrawlerURL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TaskID    string         `gorm:"not null;index" json:"task_id"` // 任务ID
	URL       string         `gorm:"size:2048;not null;index" json:"url"`
	Source    string         `gorm:"size:2048" json:"source"`      // 来源URL
	Status    int            `gorm:"default:0" json:"status"`      // 0: 未扫描, 1: 已扫描, 2: 扫描中
	Scanned   bool           `gorm:"default:false" json:"scanned"` // 是否已扫描
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 设置表名
func (CrawlerURL) TableName() string {
	return "crawler_urls"
}
