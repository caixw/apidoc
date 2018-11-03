// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/locale"
)

func newError(file, tag string, line int, msg message.Reference, vals ...interface{}) *errors.Error {
	return &errors.Error{
		File:  file,
		Line:  line,
		Field: tag,
		LocaleError: errors.LocaleError{
			MessageKey:  msg,
			MessageArgs: vals,
		},
	}
}

func (l *Lexer) SyntaxError(err *errors.Error) {
	l.h.SyntaxError(err)
}

func (l *Lexer) SyntaxWarn(err *errors.Error) {
	l.h.SyntaxWarn(err)
}

// ErrInvalidFormat 输出格式无效的错误信息
func (t *Tag) ErrInvalidFormat() {
	t.l.SyntaxError(newError(t.File, t.Name, t.Line, locale.ErrInvalidFormat))
}

// ErrDuplicateTag 输出标签重复的错误信息
func (t *Tag) ErrDuplicateTag() {
	t.l.SyntaxError(newError(t.File, t.Name, t.Line, locale.ErrDuplicateTag))
}

// ErrInvalidTag 输出无效的标签错误
func (t *Tag) ErrInvalidTag() {
	t.l.SyntaxError(newError(t.File, t.Name, t.Line, locale.ErrInvalidTag))
}
