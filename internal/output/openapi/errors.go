// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import "github.com/caixw/apidoc/internal/locale"

// Error 错误接口
type Error struct {
	Field   string
	Message string
}

func (err *Error) Error() string {
	return locale.Sprintf(locale.ErrInvalidOpenapi, err.Field, err.Message)
}

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() *Error
}

func newError(field, message string) *Error {
	return &Error{
		Field:   field,
		Message: message,
	}
}
