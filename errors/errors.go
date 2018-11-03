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

// Error 错误信息
type Error struct {
	Type  int8
	File  string
	Line  int
	Field string

	// 保存着错误内容的本地化信息。
	// 仅在返回错误信息时，才会转换成本地化内容。
	MessageKey  message.Reference
	MessageArgs []interface{}
}

// Message 输出的错误信息。
//
// 仅是错误信息，但是不会包含行号等内容。
func (err *Error) Message() string {
	return locale.Sprintf(err.MessageKey, err.MessageArgs...)
}

func (err *Error) Error() string {
	// ErrMessage = "错误信息 %s 位次于 %s:%d 的 %s",
	return locale.Sprintf(locale.ErrMessage, err.Message(), err.File, err.Line, err.Field)
}
