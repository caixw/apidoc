// SPDX-License-Identifier: MIT

package message

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// SyntaxError 表示语法错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
// apidoc 的错误基本上都是基于 SyntaxError。
type SyntaxError struct {
	locale *locale.Locale
	prev   error

	File  string
	Line  int
	Field string
}

// Message 语法错误的提示信息
func (err *SyntaxError) Message() string {
	if err.locale != nil {
		return err.locale.String()
	} else if err.prev != nil {
		return err.prev.Error()
	}

	panic("locale 与 prev 不能同时为空")
}

func (err *SyntaxError) Error() string {
	msg := err.Message()

	// ErrMessage = "%s 位次于 %s:%d 的 %s"
	return locale.Sprintf(locale.ErrMessage, msg, err.File, err.Line, err.Field)
}

// NewError 声明新的 SyntaxError 实例
func NewError(file, field string, line int, msg message.Reference, val ...interface{}) *SyntaxError {
	return &SyntaxError{
		locale: locale.NewLocale(msg, val...),
		File:   file,
		Line:   line,
		Field:  field,
	}
}

// WithError 声明 SyntaxError 实例，其中的提示信息由 err 返回
func WithError(file, field string, line int, err error) *SyntaxError {
	return &SyntaxError{
		prev:  err,
		File:  file,
		Line:  line,
		Field: field,
	}
}
