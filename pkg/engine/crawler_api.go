package engine

import (
	"github.com/projectdiscovery/gologger"
	"gorm.io/gorm"
)

// CrawlAndScan 提供给外部调用的爬虫和扫描功能
// 该函数会爬取目标URL，将爬取到的URL存入数据库，并将URL发送给扫描引擎进行漏洞检测
func CrawlAndScan(targetURL string, db *gorm.DB) error {
	// 创建爬虫选项
	opts := &CrawlerOptions{
		Depth:        3,             // 最大爬取深度
		FieldScope:   "rdn",         // 爬取范围字段
		BodyReadSize: 1024 * 1024,   // 最大响应大小
		Timeout:      10,            // 请求超时时间（秒）
		Concurrency:  10,            // 并发爬取的goroutine数量
		Parallelism:  10,            // URL处理的goroutine数量
		Delay:        0,             // 每次爬取请求之间的延迟（秒）
		RateLimit:    150,           // 每秒最大请求数
		Strategy:     "depth-first", // 访问策略（深度优先、广度优先）
	}

	// 创建扫描引擎
	scanner := NewScanner()

	// 记录开始爬取
	gologger.Info().Msgf("开始爬取目标: %s", targetURL)

	// 启动爬虫并将结果发送给扫描引擎
	err := Start(targetURL, opts, db, scanner)
	if err != nil {
		gologger.Warning().Msgf("爬取目标失败 %s: %s", targetURL, err.Error())
		return err
	}

	gologger.Info().Msgf("成功完成目标爬取: %s", targetURL)
	return nil
}
