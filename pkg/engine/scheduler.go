package engine

import (
	"context"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/seaung/yarx-go/internal/pkg/models"
	"gorm.io/gorm"
)

// Scheduler 调度器接口
type Scheduler interface {
	// Start 启动调度器
	Start(queue <-chan string, plugins []Plugin) error
	// Stop 停止调度器
	Stop()
	// AddWorker 添加工作协程
	AddWorker()
}

// SchedulerOptions 调度器配置选项
type SchedulerOptions struct {
	WorkerCount int
	DB          *gorm.DB
}

// DefaultScheduler 默认调度器实现
type DefaultScheduler struct {
	opts      *SchedulerOptions
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	workerSem chan struct{} // 工作协程信号量
}

// NewScheduler 创建一个新的调度器
func NewScheduler(opts *SchedulerOptions) Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &DefaultScheduler{
		opts:      opts,
		ctx:       ctx,
		cancel:    cancel,
		workerSem: make(chan struct{}, opts.WorkerCount),
	}
}

// Start 启动调度器
func (s *DefaultScheduler) Start(queue <-chan string, plugins []Plugin) error {
	gologger.Info().Msgf("启动调度器，工作协程数: %d", s.opts.WorkerCount)

	// 启动工作协程池
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.processQueue(queue, plugins)
	}()

	return nil
}

// Stop 停止调度器
func (s *DefaultScheduler) Stop() {
	gologger.Info().Msg("正在停止调度器...")
	s.cancel()
	s.wg.Wait()
	gologger.Info().Msg("调度器已停止")
}

// AddWorker 添加工作协程
func (s *DefaultScheduler) AddWorker() {
	// 暂未实现动态扩展工作协程
}

// processQueue 处理URL队列
func (s *DefaultScheduler) processQueue(queue <-chan string, plugins []Plugin) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case url, ok := <-queue:
			if !ok {
				return
			}

			// 获取工作协程信号量
			select {
			case s.workerSem <- struct{}{}:
				// 启动工作协程处理URL
				s.wg.Add(1)
				go func(targetURL string) {
					defer s.wg.Done()
					defer func() { <-s.workerSem }()

					s.scanURL(targetURL, plugins)
				}(url)
			case <-s.ctx.Done():
				return
			}
		}
	}
}

// scanURL 扫描URL
func (s *DefaultScheduler) scanURL(url string, plugins []Plugin) {
	gologger.Debug().Msgf("开始扫描URL: %s", url)

	// 更新URL状态为扫描中
	s.updateURLStatus(url, 1) // 1表示扫描中

	// 创建扫描上下文
	scanCtx := &ScanContext{
		URL:     url,
		Started: time.Now(),
		Results: make([]ScanResult, 0),
	}

	// 运行所有插件
	for _, plugin := range plugins {
		if plugin.Match(url) {
			result, err := plugin.Execute(scanCtx)
			if err != nil {
				gologger.Warning().Msgf("插件 %s 执行失败: %s", plugin.Name(), err.Error())
				continue
			}

			if result.Vulnerable {
				gologger.Info().Msgf("发现漏洞 [%s]: %s", plugin.Name(), url)
				// 保存漏洞结果
				s.saveVulnerability(url, plugin.Name(), result)
			}

			// 添加到扫描结果
			scanCtx.Results = append(scanCtx.Results, result)
		}
	}

	// 更新URL状态为已扫描
	s.updateURLStatus(url, 2) // 2表示已扫描
	gologger.Debug().Msgf("URL扫描完成: %s", url)
}

// updateURLStatus 更新URL扫描状态
func (s *DefaultScheduler) updateURLStatus(url string, status int) {
	// 更新数据库中的URL状态
	s.opts.DB.Model(&models.CrawlerURL{}).Where("url = ?", url).Updates(map[string]interface{}{
		"status":  status,
		"scanned": status == 2,
	})
}

// saveVulnerability 保存漏洞信息
func (s *DefaultScheduler) saveVulnerability(url, pluginName string, result ScanResult) {
	// 创建漏洞记录
	vuln := &models.Vulnerability{
		URL:         url,
		Type:        pluginName,
		Severity:    result.Severity,
		Description: result.Description,
		Payload:     result.Payload,
		Details:     result.Details,
		CreatedAt:   time.Now(),
	}

	// 保存到数据库
	s.opts.DB.Create(vuln)
}
