// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/locale"
)

func newError(file, tag string, line int, msg message.Reference, vals ...interface{}) error {
	return &errors.Error{
		File:        file,
		Line:        line,
		Field:       tag,
		MessageKey:  msg,
		MessageArgs: vals,
	}
}

// ErrInvalidFormat 返回格式无效的错误信息
func (t *Tag) ErrInvalidFormat() error {
	return newError(t.File, t.Name, t.Line, locale.ErrInvalidFormat)
}

// ErrDuplicateTag 返回标签重复的错误信息
func (t *Tag) ErrDuplicateTag() error {
	return newError(t.File, t.Name, t.Line, locale.ErrDuplicateTag)
}

// ErrInvalidTag 返回无效的标签错误
func (t *Tag) ErrInvalidTag() error {
	return newError(t.File, t.Name, t.Line, locale.ErrInvalidTag)
}
