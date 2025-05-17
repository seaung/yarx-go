package engine

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

// RedisStorage Redis存储实现，适用于分布式环境下的URL去重
type RedisStorage struct {
	DB      *gorm.DB
	redis   *redis.Client
	expire  time.Duration // URL在Redis中的过期时间
	context context.Context
}

// NewRedisStorage 创建一个新的Redis存储实例
func NewRedisStorage(db *gorm.DB, redisAddr, redisPassword string, redisDB int, expireHours int) (URLStorage, error) {
	ctx := context.Background()

	// 连接Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{
		DB:      db,
		redis:   rdb,
		expire:  time.Duration(expireHours) * time.Hour,
		context: ctx,
	}, nil
}

// SaveURL 使用Redis高效去重并将URL保存到数据库
func (s *RedisStorage) SaveURL(url, source, taskID string) error {
	// 生成Redis键
	redisKey := "url:" + taskID + ":" + url

	// 使用Redis的SETNX命令尝试设置键值，如果键已存在则返回0
	result, err := s.redis.SetNX(s.context, redisKey, 1, s.expire).Result()
	if err != nil {
		// Redis操作失败，回退到数据库检查
		var count int64
		s.DB.Model(&models.CrawlerURL{}).Where("url = ? AND task_id = ?", url, taskID).Count(&count)
		if count > 0 {
			// URL已存在于数据库中
			return nil
		}
	} else if !result {
		// 键已存在，表示URL已处理过
		return nil
	}

	// URL不存在或Redis操作失败但数据库中也不存在，添加到数据库
	crawlerURL := &models.CrawlerURL{
		URL:     url,
		Source:  source,
		TaskID:  taskID,
		Status:  0,
		Scanned: false,
	}

	return s.DB.Create(crawlerURL).Error
}

// Close 关闭Redis连接
func (s *RedisStorage) Close() error {
	return s.redis.Close()
}
