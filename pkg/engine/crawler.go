package engine

import (
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/projectdiscovery/katana/pkg/output"
	"github.com/projectdiscovery/katana/pkg/types"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

type CrawlerOptions struct {
	Depth        int
	Timeout      int
	Concurrency  int
	Delay        int
	RateLimit    int
	BodyReadSize int
	Parallelism  int
	Strategy     string
	FieldScope   string
}

func NewOptions(opts *CrawlerOptions) *types.Options {
	return &types.Options{
		MaxDepth:     opts.Depth,
		FieldScope:   opts.FieldScope,
		BodyReadSize: opts.BodyReadSize,
		Timeout:      opts.Timeout,
		Concurrency:  opts.Concurrency,
		Delay:        opts.Delay,
		Strategy:     opts.Strategy,
		Parallelism:  opts.Parallelism,
		RateLimit:    opts.RateLimit,
	}
}

// URLStorage 定义URL存储接口
type URLStorage interface {
	SaveURL(url, source, taskID string) error
}

// DBStorage 数据库存储实现
type DBStorage struct {
	DB *gorm.DB
	// 内存缓存，用于存储已处理过的URL，减少数据库查询
	urlCache map[string]bool
	mutex    sync.RWMutex
}

// NewDBStorage 创建数据库存储实例
func NewDBStorage(db *gorm.DB) URLStorage {
	return &DBStorage{
		DB:       db,
		urlCache: make(map[string]bool),
		mutex:    sync.RWMutex{},
	}
}

// SaveURL 将URL保存到数据库，如果URL已存在则不重复添加
func (s *DBStorage) SaveURL(url, source, taskID string) error {
	// 生成缓存键
	cacheKey := url + "-" + taskID

	// 首先检查内存缓存
	s.mutex.RLock()
	if _, exists := s.urlCache[cacheKey]; exists {
		s.mutex.RUnlock()
		// URL已存在于缓存中，不需要重复添加
		return nil
	}
	s.mutex.RUnlock()

	// 缓存中不存在，检查数据库
	var count int64
	s.DB.Model(&models.CrawlerURL{}).Where("url = ? AND task_id = ?", url, taskID).Count(&count)
	if count > 0 {
		// URL已存在于数据库中，添加到缓存并返回
		s.mutex.Lock()
		s.urlCache[cacheKey] = true
		s.mutex.Unlock()
		return nil
	}

	// URL不存在，添加到数据库和缓存
	crawlerURL := &models.CrawlerURL{
		URL:     url,
		Source:  source,
		TaskID:  taskID,
		Status:  0,
		Scanned: false,
	}

	err := s.DB.Create(crawlerURL).Error
	if err == nil {
		// 添加成功，更新缓存
		s.mutex.Lock()
		s.urlCache[cacheKey] = true
		s.mutex.Unlock()
	}

	return err
}

// StorageOpts URL存储选项
type StorageOpts struct {
	// 存储类型: "db"(数据库+内存缓存), "bloom"(布隆过滤器), "redis"(Redis分布式缓存)
	Type string
	// 布隆过滤器配置
	BloomExpectedItems uint    // 预期处理的URL数量
	BloomFalsePositive float64 // 可接受的误判率
	// Redis配置
	RedisAddr     string // Redis地址
	RedisPassword string // Redis密码
	RedisDB       int    // Redis数据库
	RedisExpire   int    // URL在Redis中的过期时间（小时）
}

// Start 启动爬虫并将结果存入数据库和发送给扫描引擎
func Start(target string, opts *CrawlerOptions, db *gorm.DB, scanner Scanner) error {
	return StartWithStorage(target, opts, db, scanner, nil)
}

// StartWithStorage 使用指定存储选项启动爬虫
func StartWithStorage(target string, opts *CrawlerOptions, db *gorm.DB, scanner Scanner, storageOpts *StorageOpts) error {
	// 创建存储配置
	var storageConfig *StorageConfig
	if storageOpts != nil {
		storageConfig = &StorageConfig{
			Type:               StorageType(storageOpts.Type),
			BloomExpectedItems: storageOpts.BloomExpectedItems,
			BloomFalsePositive: storageOpts.BloomFalsePositive,
			RedisAddr:          storageOpts.RedisAddr,
			RedisPassword:      storageOpts.RedisPassword,
			RedisDB:            storageOpts.RedisDB,
			RedisExpire:        storageOpts.RedisExpire,
		}
	} else {
		// 使用默认配置（布隆过滤器）
		storageConfig = DefaultStorageConfig()
	}

	// 创建URL存储实例
	storage, err := NewURLStorage(db, storageConfig)
	if err != nil {
		gologger.Error().Msgf("创建URL存储失败: %s", err.Error())
		return err
	}

	// 如果存储实现了Close接口，确保资源释放
	if closer, ok := storage.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// 生成任务ID（使用当前时间戳作为简单实现）
	taskID := time.Now().Format("20060102150405")

	// 记录爬虫任务开始
	gologger.Info().Msgf("开始爬虫任务，目标: %s，任务ID: %s", target, taskID)

	// 配置爬虫选项
	options := NewOptions(opts)

	// 设置结果回调函数
	options.OnResult = func(result output.Result) {
		url := result.Request.URL

		// 记录发现的URL
		gologger.Debug().Msgf("爬虫发现URL: %s，来源: %s，任务ID: %s", url, target, taskID)

		// 保存URL到数据库
		err := storage.SaveURL(url, target, taskID)
		if err != nil {
			gologger.Warning().Msgf("保存URL到数据库失败: %s, 错误: %s", url, err.Error())
			return
		}

		// 发送URL给扫描引擎进行漏洞检测
		go func(urlToScan string) {
			if err := scanner.Scan(urlToScan); err != nil {
				gologger.Warning().Msgf("将URL添加到扫描队列失败: %s, 错误: %s", urlToScan, err.Error())
			} else {
				gologger.Debug().Msgf("成功将URL添加到扫描队列: %s", urlToScan)
			}
		}(url)
	}

	crawlerOpts, err := types.NewCrawlerOptions(options)
	if err != nil {
		gologger.Error().Msgf("创建爬虫选项失败: %s", err.Error())
		return err
	}
	defer crawlerOpts.Close()

	crawler, err := standard.New(crawlerOpts)
	if err != nil {
		gologger.Error().Msgf("初始化爬虫失败: %s", err.Error())
		return err
	}
	defer crawler.Close()

	gologger.Info().Msgf("开始爬取目标URL: %s", target)
	err = crawler.Crawl(target)
	if err != nil {
		gologger.Error().Msgf("爬取过程中发生错误: %s", err.Error())
		return err
	}

	gologger.Info().Msgf("爬虫任务完成，目标: %s", target)
	return nil
}
