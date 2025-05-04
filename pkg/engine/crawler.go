package engine

import (
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
	SaveURL(url, source string) error
}

// DBStorage 数据库存储实现
type DBStorage struct {
	DB *gorm.DB
}

// NewDBStorage 创建数据库存储实例
func NewDBStorage(db *gorm.DB) URLStorage {
	return &DBStorage{DB: db}
}

// SaveURL 将URL保存到数据库
func (s *DBStorage) SaveURL(url, source string) error {
	crawlerURL := &models.CrawlerURL{
		URL:     url,
		Source:  source,
		Status:  0,
		Scanned: false,
	}
	return s.DB.Create(crawlerURL).Error
}

// Start 启动爬虫并将结果存入数据库和发送给扫描引擎
func Start(target string, opts *CrawlerOptions, db *gorm.DB, scanner Scanner) error {
	// 创建URL存储实例
	storage := NewDBStorage(db)

	// 配置爬虫选项
	options := NewOptions(opts)

	// 设置结果回调函数
	options.OnResult = func(result output.Result) {
		url := result.Request.URL

		// 保存URL到数据库
		err := storage.SaveURL(url, target)
		if err != nil {
			// 处理错误，可以记录日志
			return
		}

		// 发送URL给扫描引擎进行漏洞检测
		go scanner.Scan(url)
	}

	crawlerOpts, err := types.NewCrawlerOptions(options)
	if err != nil {
		return err
	}
	defer crawlerOpts.Close()

	crawler, err := standard.New(crawlerOpts)
	if err != nil {
		return err
	}
	defer crawler.Close()

	err = crawler.Crawl(target)
	if err != nil {
		return err
	}

	return nil
}
