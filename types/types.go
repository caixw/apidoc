// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package types 一些公用类型的定义
package types

import (
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/vars"
)

// Sanitizer 配置项的检测接口
type Sanitizer interface {
	Sanitize() *OptionsError
}

// OptionsError 提供对配置项错误的描述
type OptionsError struct {
	Field   string
	Message string
}

func (err *OptionsError) Error() string {
	return locale.Sprintf(locale.OptionsError, vars.ConfigFilename, err.Field, err.Message)
}
