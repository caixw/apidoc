// SPDX-License-Identifier: MIT

package message

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// LocaleError 本地化的错误接口
type LocaleError interface {
	error
	LocaleError(p *message.Printer) string
}

// SyntaxError 表示语法错误信息
//
// 无论是配置文件的错误，还是文档的语法错误，都将返回此错误。
// apidoc 的错误基本上都是基于 SyntaxError。
type SyntaxError struct {
	key  message.Reference
	args []interface{}
	prev error

	File  string
	Line  int
	Field string
}

// Message 语法错误的提示信息
func (err *SyntaxError) Message(p *message.Printer) string {
	if err.prev != nil {
		if l, ok := err.prev.(LocaleError); ok {
			return l.LocaleError(p)
		}
		return err.prev.Error()
	}

	return p.Sprintf(err.key, err.args...)
}

func (err *SyntaxError) Error() string {
	return err.LocaleError(message.NewPrinter(language.Und))
}

// LocaleError 实现 LocaleError 接口
func (err *SyntaxError) LocaleError(p *message.Printer) string {
	msg := err.Message(p)

	// TODO 根据是否有 filed，返回不同的提示内容

	// ErrMessage = "%s 位次于 %s:%d 的 %s"
	return locale.Sprintf(locale.ErrMessage, msg, err.File, err.Line, err.Field)
}

// NewError 声明新的 SyntaxError 实例
func NewError(file, field string, line int, msg message.Reference, val ...interface{}) *SyntaxError {
	return &SyntaxError{
		key:   msg,
		args:  val,
		File:  file,
		Line:  line,
		Field: field,
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
