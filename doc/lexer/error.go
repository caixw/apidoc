// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/errors"
)

// Error 输出错误信息
func (l *Lexer) Error(err *errors.Error) {
	l.h.SyntaxError(err)
}

// Warn 输出警告信息
func (l *Lexer) Warn(err *errors.Error) {
	l.h.SyntaxWarn(err)
}

// Warn 输出警告信息
func (t *Tag) Warn(key message.Reference, vals ...interface{}) {
	t.l.Warn(errors.New(t.File, t.Name, t.Line, key, vals...))
}

// Error 输出错误信息
func (t *Tag) Error(key message.Reference, vals ...interface{}) {
	t.l.Error(errors.New(t.File, t.Name, t.Line, key, vals...))
}

// WarnWithError 输出警告信息
func (t *Tag) WarnWithError(err error, key message.Reference, vals ...interface{}) {
	t.l.Warn(errors.WithError(err, t.File, t.Name, t.Line, key, vals...))
}

// ErrorWithError 输出错误信息
func (t *Tag) ErrorWithError(err error, key message.Reference, vals ...interface{}) {
	t.l.Error(errors.WithError(err, t.File, t.Name, t.Line, key, vals...))
}
