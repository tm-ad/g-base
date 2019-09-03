// Package exceptions define a interface
// for Output Specific Errors of Internal Projects and Implementation of CommonException
package exceptions

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"runtime"
)

// Exception define a custom exceptions interfaces for tmi internal projects
type Exception interface {
	// Error implementing golang's own error interface
	Error() string
	// Code gets the Code of exceptions, eg 1000001
	Code() int
	// Message gets the message of exception
	Message() string
}

// IsIException 用于判断 error 是否为自定义的 IException 接口
func IsException(err error) bool {
	_, ok := interface{}(err).(Exception)

	return ok
}

const CommonExceptionCode = 819

// CommonException 是实现Exception接口的基础异常结构
// 常规异常可通过参数的不同直接实例化此结构来支持业务错误的定义
// 更高级的归类可嵌入结构体实现
type CommonException struct {
	code    int
	message string
}

// expJson 用于输出JSON
type expJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 将异常信息输出为JSON字符串，以用于传输后的进一步处理
func (e *CommonException) Error() string {
	ej := expJson{
		Code:    e.Code(),
		Message: e.Message(),
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	j, err := json.MarshalToString(ej)
	if err == nil {
		return j
	} else {
		// fall back to string
		return fmt.Sprintf("%s (%d)", e.message, e.code)
	}
}

// Code 获取异常所对应的错误业务异常Code
func (e *CommonException) Code() int {
	return e.code
}

// Message 获取异常所对应的错误业务异常消息
func (e *CommonException) Message() string {
	return e.message
}

func newCommonException(code int, message string, causes ...error) Exception {
	return &CommonException{
		code:    code,
		message: message,
	}
}

// New 创建一个Code为819的常规异常，如需要自定义code，请使用 exceptions.Code
func New(message string) Exception {
	return newCommonException(CommonExceptionCode, message)
}

// NewF 创建一个Code为819的常规异常，如需要自定义code，请使用 exceptions.Code
func NewF(format string, args ...interface{}) Exception {
	return New(fmt.Sprintf(format, args...))
}

// Code 创建一个自定义Code的异常
func Code(code int, message string) Exception {
	return newCommonException(code, message)
}

// CodeF 创建一个自定义Code的异常
func CodeF(code int, format string, args ...interface{}) Exception {
	return Code(code, fmt.Sprintf(format, args...))
}

// CallStack 获取调用堆栈
func CallStack() []string {
	stack := []string{}

	for i := 1; ; i++ {
		_, file, line, got := runtime.Caller(i)
		if !got {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}

	return stack
}
