package engine

import (
	"sync"

	"github.com/projectdiscovery/gologger"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

// BloomStorage 基于布隆过滤器的URL存储实现
type BloomStorage struct {
	DB          *gorm.DB
	bloomFilter *BloomFilter
	urlCache    map[string]bool // 用于确认性检查的小型缓存
	mutex       sync.RWMutex
}

// NewBloomStorage 创建一个新的布隆过滤器存储实例
// expectedItems: 预期处理的URL数量
// falsePositiveRate: 可接受的误判率 (0.0 - 1.0)
func NewBloomStorage(db *gorm.DB, expectedItems uint, falsePositiveRate float64) URLStorage {
	return &BloomStorage{
		DB:          db,
		bloomFilter: NewBloomFilter(expectedItems, falsePositiveRate),
		urlCache:    make(map[string]bool),
		mutex:       sync.RWMutex{},
	}
}

// SaveURL 使用布隆过滤器高效去重并将URL保存到数据库
func (s *BloomStorage) SaveURL(url, source, taskID string) error {
	// 生成缓存键
	cacheKey := url + "-" + taskID

	// 首先检查布隆过滤器
	if s.bloomFilter.Contains(cacheKey) {
		// 布隆过滤器表明URL可能存在，进一步检查确认缓存
		s.mutex.RLock()
		if _, exists := s.urlCache[cacheKey]; exists {
			s.mutex.RUnlock()
			// URL确认存在于缓存中，不需要重复添加
			return nil
		}
		s.mutex.RUnlock()

		// 布隆过滤器可能有误判，检查数据库确认
		var count int64
		s.DB.Model(&models.CrawlerURL{}).Where("url = ? AND task_id = ?", url, taskID).Count(&count)
		if count > 0 {
			// URL确实存在于数据库中，添加到确认缓存
			s.mutex.Lock()
			s.urlCache[cacheKey] = true
			s.mutex.Unlock()
			return nil
		}
	}

	// URL不存在或布隆过滤器未命中，添加到数据库
	crawlerURL := &models.CrawlerURL{
		URL:     url,
		Source:  source,
		TaskID:  taskID,
		Status:  0,
		Scanned: false,
	}

	err := s.DB.Create(crawlerURL).Error
	if err == nil {
		// 添加成功，更新布隆过滤器和确认缓存
		s.bloomFilter.Add(cacheKey)

		// 只在确认缓存中保存一部分最近的URL，防止内存占用过大
		s.mutex.Lock()
		// 如果缓存过大，可以考虑清理部分旧数据
		if len(s.urlCache) > 10000 { // 设置一个合理的缓存大小上限
			gologger.Debug().Msgf("确认缓存达到上限，清理缓存")
			s.urlCache = make(map[string]bool) // 简单实现：直接清空缓存
		}
		s.urlCache[cacheKey] = true
		s.mutex.Unlock()
	}

	return err
}
