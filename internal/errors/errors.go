// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package errors 公用的错误信息
package errors

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/locale"
)

// Error 错误信息
type Error struct {
	File  string
	Line  int
	Field string

	// 保存着错误内容的本地化信息。
	// 仅在返回错误信息时，才会转换成本地化内容。
	MessageKey  message.Reference
	MessageArgs []interface{}
}

func (err *Error) Error() string {
	msg := locale.Sprintf(err.MessageKey, err.MessageArgs...)
	return locale.Sprintf(locale.ErrSyntax, err.File, err.Line, msg)
}
