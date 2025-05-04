package models

import "gorm.io/gorm"

// AutoMigrate 自动迁移数据库模型
func AutoMigrate(db *gorm.DB) error {
	// 在这里添加所有需要自动迁移的模型
	return db.AutoMigrate(
		&CrawlerURL{},
	)
}
