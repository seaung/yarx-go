package errno

import "fmt"

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

// 实现error接口
func (e *Errno) Error() string {
	return e.Message
}

// 设置错误消息
func (e *Errno) SetMessage(format string, args ...interface{}) *Errno {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// 从err中解析业务错误码和错误信息
func Decode(err error) (int, string, string) {
	if err == nil {
		return HTTPOK.HTTP, HTTPOK.Code, HTTPOK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	return InternalServerError.HTTP, InternalServerError.Code, err.Error()
}
