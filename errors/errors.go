// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package errors 公用的错误信息
package errors

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/locale"
)

// 错误分类
const (
	SyntaxError int8 = iota + 1
	SyntaxWarn
	Other
)

// LocaleError 本地化的错误信息
type LocaleError struct {
	MessageKey  message.Reference
	MessageArgs []interface{}
}

// Error 错误信息
type Error struct {
	LocaleError
	Message string

	Type  int8
	File  string
	Line  int
	Field string
}

func (err *LocaleError) Error() string {
	return locale.Sprintf(err.MessageKey, err.MessageArgs...)
}

func (err *Error) Error() string {
	msg := err.Message
	if msg == "" {
		msg = err.LocaleError.Error()
	}

	return locale.Sprintf(locale.ErrMessage, msg, err.File, err.Line, err.Field)
}

// New 声明新的 Error 实例
func New(file, field string, line int, msg message.Reference, vals ...interface{}) *Error {
	return &Error{
		File:  file,
		Line:  line,
		Field: field,
		LocaleError: LocaleError{
			MessageKey:  msg,
			MessageArgs: vals,
		},
	}
}
