package engine

import (
	"errors"
	"io/fs"
	"path/filepath"
	"plugin"
	"time"

	"github.com/projectdiscovery/gologger"
)

// ScanContext 扫描上下文
type ScanContext struct {
	URL     string
	Started time.Time
	Results []ScanResult
}

// ScanResult 扫描结果
type ScanResult struct {
	Vulnerable  bool
	Severity    string // 严重程度：high, medium, low, info
	Description string
	Payload     string
	Details     string
}

// Plugin 漏洞检测插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string
	// Description 返回插件描述
	Description() string
	// Match 判断URL是否适用于该插件
	Match(url string) bool
	// Execute 执行漏洞检测
	Execute(ctx *ScanContext) (ScanResult, error)
}

// BasePlugin 基础插件实现
type BasePlugin struct {
	name        string
	description string
}

// Name 返回插件名称
func (p *BasePlugin) Name() string {
	return p.name
}

// Description 返回插件描述
func (p *BasePlugin) Description() string {
	return p.description
}

// Match 默认匹配所有URL
func (p *BasePlugin) Match(url string) bool {
	return true
}

// Execute 默认执行方法，需要被子类覆盖
func (p *BasePlugin) Execute(ctx *ScanContext) (ScanResult, error) {
	return ScanResult{}, errors.New("未实现的方法")
}

// LoadPlugins 从指定目录加载插件
func LoadPlugins(pluginDir string) []Plugin {
	plugins := make([]Plugin, 0)

	// 如果插件目录为空，返回内置插件
	if pluginDir == "" {
		return loadBuiltinPlugins()
	}

	// 遍历插件目录
	err := filepath.WalkDir(pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理.so文件
		if !d.IsDir() && filepath.Ext(path) == ".so" {
			// 加载插件
			p, err := plugin.Open(path)
			if err != nil {
				gologger.Warning().Msgf("加载插件失败 %s: %s", path, err.Error())
				return nil
			}

			// 查找插件符号
			sym, err := p.Lookup("Plugin")
			if err != nil {
				gologger.Warning().Msgf("插件符号查找失败 %s: %s", path, err.Error())
				return nil
			}

			// 类型断言
			plg, ok := sym.(Plugin)
			if !ok {
				gologger.Warning().Msgf("插件类型断言失败 %s", path)
				return nil
			}

			// 添加到插件列表
			plugins = append(plugins, plg)
			gologger.Info().Msgf("成功加载插件: %s - %s", plg.Name(), plg.Description())
		}

		return nil
	})

	if err != nil {
		gologger.Warning().Msgf("遍历插件目录失败: %s", err.Error())
	}

	// 如果没有找到外部插件，加载内置插件
	if len(plugins) == 0 {
		return loadBuiltinPlugins()
	}

	return plugins
}

// loadBuiltinPlugins 加载内置插件
func loadBuiltinPlugins() []Plugin {
	plugins := make([]Plugin, 0)

	// 添加SQL注入检测插件
	plugins = append(plugins, &SQLInjectionPlugin{
		BasePlugin: BasePlugin{
			name:        "sql-injection",
			description: "SQL注入漏洞检测",
		},
	})

	// 添加XSS检测插件
	plugins = append(plugins, &XSSPlugin{
		BasePlugin: BasePlugin{
			name:        "xss",
			description: "跨站脚本攻击检测",
		},
	})

	// 添加目录遍历检测插件
	plugins = append(plugins, &DirectoryTraversalPlugin{
		BasePlugin: BasePlugin{
			name:        "directory-traversal",
			description: "目录遍历漏洞检测",
		},
	})

	gologger.Info().Msgf("已加载 %d 个内置插件", len(plugins))
	return plugins
}

// SQLInjectionPlugin SQL注入检测插件
type SQLInjectionPlugin struct {
	BasePlugin
}

// Match 判断URL是否适用于SQL注入检测
func (p *SQLInjectionPlugin) Match(url string) bool {
	// 检查URL是否包含参数
	return true
}

// Execute 执行SQL注入检测
func (p *SQLInjectionPlugin) Execute(ctx *ScanContext) (ScanResult, error) {
	// 这里实现SQL注入检测逻辑
	// 实际应用中应该发送带有SQL注入payload的请求并分析响应

	// 模拟检测结果
	return ScanResult{
		Vulnerable:  false,
		Severity:    "high",
		Description: "SQL注入漏洞",
		Payload:     "' OR 1=1 --",
		Details:     "未检测到SQL注入漏洞",
	}, nil
}

// XSSPlugin 跨站脚本攻击检测插件
type XSSPlugin struct {
	BasePlugin
}

// Match 判断URL是否适用于XSS检测
func (p *XSSPlugin) Match(url string) bool {
	// 检查URL是否包含参数
	return true
}

// Execute 执行XSS检测
func (p *XSSPlugin) Execute(ctx *ScanContext) (ScanResult, error) {
	// 这里实现XSS检测逻辑
	// 实际应用中应该发送带有XSS payload的请求并分析响应

	// 模拟检测结果
	return ScanResult{
		Vulnerable:  false,
		Severity:    "medium",
		Description: "跨站脚本攻击漏洞",
		Payload:     "<script>alert(1)</script>",
		Details:     "未检测到XSS漏洞",
	}, nil
}

// DirectoryTraversalPlugin 目录遍历检测插件
type DirectoryTraversalPlugin struct {
	BasePlugin
}

// Match 判断URL是否适用于目录遍历检测
func (p *DirectoryTraversalPlugin) Match(url string) bool {
	// 检查URL是否包含文件路径参数
	return true
}

// Execute 执行目录遍历检测
func (p *DirectoryTraversalPlugin) Execute(ctx *ScanContext) (ScanResult, error) {
	// 这里实现目录遍历检测逻辑
	// 实际应用中应该发送带有目录遍历payload的请求并分析响应

	// 模拟检测结果
	return ScanResult{
		Vulnerable:  false,
		Severity:    "medium",
		Description: "目录遍历漏洞",
		Payload:     "../../../etc/passwd",
		Details:     "未检测到目录遍历漏洞",
	}, nil
}
