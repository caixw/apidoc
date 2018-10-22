// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package options 配置项的一些处理操作。
//
// 被 input 和 output 调用。
package options

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/locale"
)

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() error
}

// FieldError 错误接口
type FieldError struct {
	Field       string
	MessageKey  message.Reference
	MessageArgs []interface{}
}

func (err *FieldError) Error() string {
	msg := locale.Sprintf(err.MessageKey, err.MessageArgs...)
	return locale.Sprintf(locale.ErrOptions, err.Field, msg)
}

// NewFieldError 声明一个新的错误对象
func NewFieldError(field string, key message.Reference, args ...interface{}) *FieldError {
	return &FieldError{
		Field:       field,
		MessageKey:  key,
		MessageArgs: args,
	}
}
