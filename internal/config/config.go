// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package config 配置文件的错误处理内容
package config

import (
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/vars"
)

// Error 错误接口
type Error struct {
	Field   string
	Message string
}

func (err *Error) Error() string {
	return locale.Sprintf(locale.ErrConfig, vars.ConfigFilename, err.Field, err.Message)
}

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() error
}

// New 声明一个新的错误对象
func New(field, message string) *Error {
	return &Error{
		Field:   field,
		Message: message,
	}
}
