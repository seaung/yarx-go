package engine

import (
	"net/http"
	"time"

	"github.com/projectdiscovery/gologger"
)

// Scanner 定义扫描引擎接口
type Scanner interface {
	// Scan 扫描URL查找漏洞
	Scan(url string) error
}

// DefaultScanner 默认扫描引擎实现
type DefaultScanner struct {
	client    *http.Client
	UserAgent string
}

// ScannerOptions 扫描器配置选项
type ScannerOptions struct {
	Timeout   int
	UserAgent string
}

// NewScanner 创建一个新的扫描引擎实例
func NewScanner() Scanner {
	return &DefaultScanner{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
}

// NewScannerWithOptions 使用自定义选项创建扫描引擎实例
func NewScannerWithOptions(opts *ScannerOptions) Scanner {
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	if opts.UserAgent != "" {
		userAgent = opts.UserAgent
	}

	timeout := 10
	if opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	return &DefaultScanner{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		UserAgent: userAgent,
	}
}

// Scan 实现Scanner接口的Scan方法
func (s *DefaultScanner) Scan(url string) error {
	gologger.Debug().Msgf("扫描URL: %s", url)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", s.UserAgent)

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 记录响应状态码
	gologger.Debug().Msgf("URL %s 响应状态码: %d", url, resp.StatusCode)

	return nil
}
