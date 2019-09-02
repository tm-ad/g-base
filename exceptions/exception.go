// Package exceptions define a interface
// for Output Specific Errors of Internal Projects and Implementation of CommonException
package exceptions

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// Exception define a custom exceptions interfaces for tmi internal projects
type Exception interface {
	// Error implementing golang's own error interface
	Error() string
	// Code gets the Code of exceptions, eg 1000001
	Code() int
	// Message gets the message of exception
	Message() string
	// Cause gets the error stack
	Causes() []error
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
	causes  []error
}

// expJson 用于输出JSON
type expJson struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Causes  []string `json:"causes"`
}

// Error 将异常信息输出为JSON字符串，以用于传输后的进一步处理
func (e *CommonException) Error() string {
	ej := expJson{
		Code:    e.Code(),
		Message: e.Message(),
		Causes:  []string{},
	}
	cs := e.Causes()
	if cs != nil && len(cs) > 0 {
		for i := 0; i < len(cs); i++ {
			ej.Causes = append(ej.Causes, cs[i].Error())
		}
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

// Cause 获取已发此异常的原因（即原始异常）
func (e *CommonException) Causes() []error {
	return e.causes
}

func newCommonException(code int, message string, causes ...error) Exception {
	return &CommonException{
		code:    code,
		message: message,
		causes:  causes,
	}
}

// New 创建一个Code为819的常规异常，如需要自定义code，请使用 exceptions.Code
func New(message string, causes ...error) Exception {
	return newCommonException(CommonExceptionCode, message, causes...)
}

// Code 创建一个自定义Code的异常
func Code(code int, message string, causes ...error) Exception {
	return newCommonException(code, message, causes...)
}
