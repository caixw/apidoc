// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/errors"
)

func (l *lexer) err(err *errors.Error) {
	l.h.SyntaxError(err)
}

func (l *lexer) warn(err *errors.Error) {
	l.h.SyntaxWarn(err)
}

func (t *lexerTag) warn(key message.Reference, vals ...interface{}) {
	t.l.warn(errors.New(t.File, t.Name, t.Line, key, vals...))
}

func (t *lexerTag) err(key message.Reference, vals ...interface{}) {
	t.l.err(errors.New(t.File, t.Name, t.Line, key, vals...))
}

func (t *lexerTag) warnWithError(err error, key message.Reference, vals ...interface{}) {
	t.l.warn(errors.WithError(err, t.File, t.Name, t.Line, key, vals...))
}

func (t *lexerTag) errWithError(err error, key message.Reference, vals ...interface{}) {
	t.l.err(errors.WithError(err, t.File, t.Name, t.Line, key, vals...))
}

func (api *API) err(tag string, key message.Reference, vals ...interface{}) *errors.Error {
	return errors.New(api.file, tag, api.line, key, vals...)
}
