package engine

// Scanner 定义扫描引擎接口
type Scanner interface {
	// Scan 扫描URL查找漏洞
	Scan(url string) error
}

// DefaultScanner 默认扫描引擎实现
type DefaultScanner struct{}

// NewScanner 创建一个新的扫描引擎实例
func NewScanner() Scanner {
	return &DefaultScanner{}
}

// Scan 实现Scanner接口的Scan方法
func (s *DefaultScanner) Scan(url string) error {
	return nil
}
