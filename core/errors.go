// SPDX-License-Identifier: MIT

package core

import (
	"net/http"

	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// HTTPError 表示 HTTP 状态码的错误
type HTTPError struct {
	Code    int
	Message string
}

// SyntaxError 表示语法错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
// apidoc 的错误基本上都是基于 SyntaxError。
type SyntaxError struct {
	Location Location
	Message  string
	Field    string
}

func (err HTTPError) Error() string {
	return err.Message
}

// NewHTTPError 声明 HTTPError 实例
func NewHTTPError(code int, msg string) *HTTPError {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return &HTTPError{
		Code:    code,
		Message: msg,
	}
}

func (err *SyntaxError) Error() string {
	detail := err.Location.String()

	if err.Field != "" {
		detail += ":" + err.Field
	}

	// ErrMessage = "%s 位次于 %s:%d"
	return locale.Sprintf(locale.ErrMessage, err.Message, detail)
}

// NewLocaleError 本地化的错误信息
//
// 其中的 msg 和 val 会被转换成本地化的内容保存。
func NewLocaleError(loc Location, field string, msg message.Reference, val ...interface{}) *SyntaxError {
	return &SyntaxError{
		Message:  locale.Sprintf(msg, val...),
		Location: loc,
		Field:    field,
	}
}

// WithError 声明 SyntaxError 实例，其中的提示信息由 err 返回
func WithError(loc Location, field string, err error) *SyntaxError {
	return &SyntaxError{
		Message:  err.Error(),
		Location: loc,
		Field:    field,
	}
}
