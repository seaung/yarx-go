package engine

import (
	"context"
	"sync"

	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/projectdiscovery/gologger"
	"gorm.io/gorm"
)

// Engine 定义漏洞扫描引擎接口
type Engine interface {
	// Start 启动引擎
	Start() error
	// Stop 停止引擎
	Stop() error
	// AddTask 添加扫描任务
	AddTask(target string) error
}

// DefaultEngine 默认引擎实现
type DefaultEngine struct {
	db        *gorm.DB
	scheduler Scheduler
	crawler   *CrawlerOptions
	scanner   Scanner
	plugins   []Plugin
	queue     chan string // URL扫描队列
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	running   bool
	mutex     sync.Mutex
}

// EngineOptions 引擎配置选项
type EngineOptions struct {
	DB           *gorm.DB
	CrawlerOpts  *CrawlerOptions
	QueueSize    int
	WorkerCount  int
	PluginFolder string
}

// NewEngine 创建一个新的扫描引擎
func NewEngine(opts *EngineOptions) Engine {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建扫描器
	scanner := NewScanner()

	// 创建调度器
	scheduler := NewScheduler(&SchedulerOptions{
		WorkerCount: opts.WorkerCount,
		DB:          opts.DB,
	})

	// 加载插件
	plugins := LoadPlugins(opts.PluginFolder)

	return &DefaultEngine{
		db:        opts.DB,
		scheduler: scheduler,
		crawler:   opts.CrawlerOpts,
		scanner:   scanner,
		plugins:   plugins,
		queue:     make(chan string, opts.QueueSize),
		ctx:       ctx,
		cancel:    cancel,
		running:   false,
	}
}

// Start 启动引擎
func (e *DefaultEngine) Start() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.running {
		return nil
	}

	// 启动调度器
	err := e.scheduler.Start(e.queue, e.plugins)
	if err != nil {
		return err
	}

	e.running = true
	gologger.Info().Msg("漏洞扫描引擎已启动")
	return nil
}

// Stop 停止引擎
func (e *DefaultEngine) Stop() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !e.running {
		return nil
	}

	// 取消上下文
	e.cancel()

	// 停止调度器
	e.scheduler.Stop()

	// 等待所有任务完成
	e.wg.Wait()

	e.running = false
	gologger.Info().Msg("漏洞扫描引擎已停止")
	return nil
}

// AddTask 添加扫描任务
func (e *DefaultEngine) AddTask(target string) error {
	// 启动爬虫爬取目标URL
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		gologger.Info().Msgf("开始爬取目标: %s", target)
		err := Start(target, e.crawler, e.db, e)
		if err != nil {
			gologger.Warning().Msgf("爬取目标失败 %s: %s", target, err.Error())
			return
		}
		gologger.Info().Msgf("成功完成目标爬取: %s", target)
	}()

	return nil
}

// ProcessURL 实现Scanner接口，将URL添加到扫描队列
func (e *DefaultEngine) Scan(url string) error {
	select {
	case <-e.ctx.Done():
		return nil
	case e.queue <- url:
		gologger.Debug().Msgf("URL已添加到扫描队列: %s", url)
	}
	return nil
}

// HandleMachineryTask 处理来自Machinery消息队列的任务
func HandleMachineryTask(db *gorm.DB, engine Engine) func(*tasks.Signature) error {
	return func(signature *tasks.Signature) error {
		// 从任务参数中获取目标URL
		if len(signature.Args) == 0 {
			return nil
		}

		target, ok := signature.Args[0].Value.(string)
		if !ok {
			return nil
		}

		// 添加到扫描引擎
		return engine.AddTask(target)
	}
}
