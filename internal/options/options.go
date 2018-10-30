// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package options 配置项的一些处理操作。
//
// 被 input 和 output 调用。
package options

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() error
}
