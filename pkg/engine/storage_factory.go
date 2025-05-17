package engine

import (
	"errors"

	"github.com/projectdiscovery/gologger"
	"gorm.io/gorm"
)

// StorageType 存储类型枚举
type StorageType string

const (
	// StorageTypeDB 数据库存储（带内存缓存）
	StorageTypeDB StorageType = "db"
	// StorageTypeBloom 布隆过滤器存储
	StorageTypeBloom StorageType = "bloom"
	// StorageTypeRedis Redis存储
	StorageTypeRedis StorageType = "redis"
)

// StorageConfig 存储配置
type StorageConfig struct {
	// 存储类型
	Type StorageType

	// 布隆过滤器配置
	BloomExpectedItems uint    // 预期处理的URL数量
	BloomFalsePositive float64 // 可接受的误判率

	// Redis配置
	RedisAddr     string // Redis地址
	RedisPassword string // Redis密码
	RedisDB       int    // Redis数据库
	RedisExpire   int    // URL在Redis中的过期时间（小时）
}

// DefaultStorageConfig 返回默认存储配置
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		Type:               StorageTypeBloom, // 默认使用布隆过滤器
		BloomExpectedItems: 100000,           // 默认预期10万个URL
		BloomFalsePositive: 0.001,            // 默认误判率0.1%
		RedisAddr:          "localhost:6379",
		RedisPassword:      "",
		RedisDB:            0,
		RedisExpire:        24, // 默认24小时
	}
}

// NewURLStorage 创建URL存储实例
func NewURLStorage(db *gorm.DB, config *StorageConfig) (URLStorage, error) {
	if config == nil {
		config = DefaultStorageConfig()
	}

	switch config.Type {
	case StorageTypeDB:
		gologger.Info().Msgf("使用数据库存储（带内存缓存）进行URL去重")
		return NewDBStorage(db), nil

	case StorageTypeBloom:
		gologger.Info().Msgf("使用布隆过滤器存储进行URL去重（预期URL数：%d，误判率：%.4f）",
			config.BloomExpectedItems, config.BloomFalsePositive)
		return NewBloomStorage(db, config.BloomExpectedItems, config.BloomFalsePositive), nil

	case StorageTypeRedis:
		gologger.Info().Msgf("使用Redis存储进行URL去重（地址：%s）", config.RedisAddr)
		return NewRedisStorage(db, config.RedisAddr, config.RedisPassword,
			config.RedisDB, config.RedisExpire)

	default:
		return nil, errors.New("不支持的存储类型")
	}
}
