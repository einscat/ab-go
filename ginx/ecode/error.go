package ecode

import (
	"fmt"
)

// Error 标准业务错误结构
type Error struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Details any    `json:"details,omitempty"`
	cause   error  `json:"-"`
}

// New 构造函数
func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// Unwrap 支持 errors.Is/As
func (e *Error) Unwrap() error {
	return e.cause
}

// WithDetails 链式调用：添加详情 (返回新实例，保证不可变性)
func (e *Error) WithDetails(details any) *Error {
	newErr := *e
	newErr.Details = details
	return &newErr
}

// WithCause 链式调用：记录底层错误 (用于日志，不返回给前端)
func (e *Error) WithCause(err error) *Error {
	newErr := *e
	newErr.cause = err
	return &newErr
}

// WithMsg 链式调用：临时修改提示语
func (e *Error) WithMsg(msg string) *Error {
	newErr := *e
	newErr.Msg = msg
	return &newErr
}
