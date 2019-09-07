// SPDX-License-Identifier: MIT

// Package errors 公用的错误信息
package errors

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// 错误分类
const (
	SyntaxError int8 = iota + 1
	SyntaxWarn
	Other
)

// LocaleError 本地化的错误信息
type LocaleError struct {
	Key  message.Reference
	Args []interface{}
}

// Error 错误信息
type Error struct {
	LocaleError
	prev error // 前一条错误信息

	Type  int8
	File  string
	Line  int
	Field string
}

func (err *LocaleError) Error() string {
	return locale.Sprintf(err.Key, err.Args...)
}

func (err *Error) Error() string {
	msg := err.LocaleError.Error()

	if err.prev == nil {
		// ErrMessage = "%s 位次于 %s:%d 的 %s"
		return locale.Sprintf(locale.ErrMessage, msg, err.File, err.Line, err.Field)
	}

	// ErrMessageWithError = "%s[%s] 位于 %s:%d 的 %s"
	return locale.Sprintf(locale.ErrMessageWithError, msg, err.prev.Error(), err.File, err.Line, err.Field)
}

// NewLocaleError 声明 LocaleError
func NewLocaleError(key message.Reference, val ...interface{}) *LocaleError {
	return &LocaleError{
		Key:  key,
		Args: val,
	}
}

// New 声明新的 Error 实例
func New(file, field string, line int, msg message.Reference, vals ...interface{}) *Error {
	return &Error{
		File:  file,
		Line:  line,
		Field: field,
		LocaleError: LocaleError{
			Key:  msg,
			Args: vals,
		},
	}
}

// WithError 返回一条带错误内容的 Error 实例
func WithError(err error, file, field string, line int, msg message.Reference, vals ...interface{}) *Error {
	return &Error{
		prev:  err,
		File:  file,
		Line:  line,
		Field: field,
		LocaleError: LocaleError{
			Key:  msg,
			Args: vals,
		},
	}
}
