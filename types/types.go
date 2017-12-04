// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package types 一些公用类型的定义
package types

// Sanitizer 配置项的检测接口
type Sanitizer interface {
	Sanitize() *OptionsError
}
