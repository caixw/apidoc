// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import "fmt"

// Error 错误接口
type Error struct {
	Field   string
	Message string
}

func (err *Error) Error() string {
	return fmt.Sprintf("在字段 %s 上发生以下错误：%s", err.Field, err.Message)
}

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() Error
}

func newError(field, message string) *Error {
	return &Error{
		Field:   field,
		Message: message,
	}
}
