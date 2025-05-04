package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLOptions struct {
	Host    string
	Name    string
	Pass    string
	Port    int
	DB      string
	Level   int
	MaxOpen int
	MaxIdle int
	MaxLife time.Duration
}

func (s *SQLOptions) getDsn() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		s.Name,
		s.Pass,
		s.Host,
		s.DB,
		true,
		"Local")
}

func NewConnection(opts *SQLOptions) (*gorm.DB, error) {
	level := logger.Silent
	if opts.Level != 0 {
		level = logger.LogLevel(opts.Level)
	}

	db, err := gorm.Open(postgres.Open(opts.getDsn()), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// SetMaxOpenConns 设置到数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpen)

	// SetConnMaxLifetime 设置连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(opts.MaxLife)

	// SetMaxIdleConns 设置空闲连接池的最大连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdle)
	return db, nil
}
