// SPDX-License-Identifier: MIT

package core

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// HTTPError 表示 HTTP 状态码的错误
type HTTPError struct {
	locale.Err
	Code int
}

// SyntaxError 表示语法错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
// apidoc 的错误基本上都是基于 SyntaxError。
type SyntaxError struct {
	Err      error
	Location Location
	Field    string
}

// NewHTTPError 声明 HTTPError 实例
func NewHTTPError(code int, key message.Reference, v ...interface{}) *HTTPError {
	return &HTTPError{
		Err:  locale.Err{Key: key, Values: v},
		Code: code,
	}
}

func (err *SyntaxError) Error() string {
	detail := err.Location.String()

	if err.Field != "" {
		detail += ":" + err.Field
	}

	// ErrMessage = "%s 位次于 %s:%d"
	return locale.Sprintf(locale.ErrMessage, err.Err.Error(), detail)
}

// Unwrap 实现 errors.Unwrap 接口
func (err *SyntaxError) Unwrap() error {
	return err.Err
}

// Is 实现 errors.Is 接口
func (err *SyntaxError) Is(target error) bool {
	return err.Err == target
}

// NewSyntaxError 本地化的错误信息
//
// 其中的 msg 和 val 会被转换成本地化的内容保存。
func NewSyntaxError(loc Location, field string, key message.Reference, val ...interface{}) *SyntaxError {
	return &SyntaxError{
		Err:      locale.NewError(key, val...),
		Location: loc,
		Field:    field,
	}
}

// NewSyntaxErrorWithError 声明 SyntaxError 实例，其中的提示信息由 err 返回
func NewSyntaxErrorWithError(loc Location, field string, err error) *SyntaxError {
	if serr, ok := err.(*SyntaxError); ok {
		err = serr.Err
	}
	return &SyntaxError{
		Err:      err,
		Location: loc,
		Field:    field,
	}
}
