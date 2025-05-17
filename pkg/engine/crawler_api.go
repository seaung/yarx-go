package engine

import (
	"github.com/projectdiscovery/gologger"
	"gorm.io/gorm"
)

// CrawlAndScan 提供给外部调用的爬虫和扫描功能
// 该函数会爬取目标URL，将爬取到的URL存入数据库，并将URL发送给扫描引擎进行漏洞检测
// CrawlOptions 爬虫配置选项
type CrawlOptions struct {
	// 爬虫基本选项
	CrawlerOpts *CrawlerOptions
	// URL存储选项
	StorageOpts *StorageOpts
}

// DefaultCrawlOptions 返回默认爬虫配置
func DefaultCrawlOptions() *CrawlOptions {
	return &CrawlOptions{
		CrawlerOpts: &CrawlerOptions{
			Depth:        3,             // 最大爬取深度
			FieldScope:   "rdn",         // 爬取范围字段
			BodyReadSize: 1024 * 1024,   // 最大响应大小
			Timeout:      10,            // 请求超时时间（秒）
			Concurrency:  10,            // 并发爬取的goroutine数量
			Parallelism:  10,            // URL处理的goroutine数量
			Delay:        0,             // 每次爬取请求之间的延迟（秒）
			RateLimit:    150,           // 每秒最大请求数
			Strategy:     "depth-first", // 访问策略（深度优先、广度优先）
		},
		// 默认使用布隆过滤器存储
		StorageOpts: &StorageOpts{
			Type:               "bloom",
			BloomExpectedItems: 100000,
			BloomFalsePositive: 0.001,
		},
	}
}

// CrawlAndScan 使用默认配置爬取目标URL并进行扫描
func CrawlAndScan(targetURL string, db *gorm.DB) error {
	return CrawlAndScanWithOptions(targetURL, db, nil)
}

// CrawlAndScanWithOptions 使用自定义配置爬取目标URL并进行扫描
func CrawlAndScanWithOptions(targetURL string, db *gorm.DB, options *CrawlOptions) error {
	// 使用默认配置或自定义配置
	if options == nil {
		options = DefaultCrawlOptions()
	}

	// 创建扫描引擎
	scanner := NewScanner()

	// 记录开始爬取
	gologger.Info().Msgf("开始爬取目标: %s", targetURL)

	// 启动爬虫并将结果发送给扫描引擎
	// StartWithStorage函数会将爬取到的URL存入数据库并送入扫描队列
	err := StartWithStorage(targetURL, options.CrawlerOpts, db, scanner, options.StorageOpts)
	if err != nil {
		gologger.Warning().Msgf("爬取目标失败 %s: %s", targetURL, err.Error())
		return err
	}

	gologger.Info().Msgf("成功完成目标爬取: %s", targetURL)
	return nil
}
