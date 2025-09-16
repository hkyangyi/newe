package utils

import (
	"fmt"
	"runtime"
)

// AppError 自定义应用错误
type AppError struct {
	Code    int    // 错误代码
	Message string // 错误信息
	File    string // 发生错误的文件
	Line    int    // 发生错误的行号
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("错误代码: %d, 消息: %s, 位置: %s:%d, 原始错误: %v",
			e.Code, e.Message, e.File, e.Line, e.Err)
	}
	return fmt.Sprintf("错误代码: %d, 消息: %s, 位置: %s:%d",
		e.Code, e.Message, e.File, e.Line)
}

// NewError 创建新的应用错误
func NewError(code int, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code:    code,
		Message: message,
		File:    file,
		Line:    line,
		Err:     err,
	}
}

// WrapError 包装现有错误
func WrapError(code int, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code:    code,
		Message: message,
		File:    file,
		Line:    line,
		Err:     err,
	}
}

// 常用错误代码
const (
	ErrCodeDatabase     = 1001 // 数据库错误
	ErrCodeValidation   = 1002 // 验证错误
	ErrCodeAuth         = 1003 // 认证错误
	ErrCodePermission   = 1004 // 权限错误
	ErrCodeSystem       = 1005 // 系统错误
	ErrCodeNetwork      = 1006 // 网络错误
	ErrCodeConfig       = 1007 // 配置错误
	ErrCodeBusiness     = 1008 // 业务错误
	ErrCodeNotFound     = 1009 // 未找到
	ErrCodeTimeout      = 1010 // 超时错误
)