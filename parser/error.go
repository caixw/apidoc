// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/caixw/apidoc/locale"
	"golang.org/x/text/message"
)

const (
	typeWarn int8 = iota + 1
	typeError
)

// 语法错误类型
type syntaxError struct {
	File        string
	Line        int
	Type        int8
	MessageKey  message.Reference
	MessageArgs []interface{}
}

func (err *syntaxError) Error() string {
	msg := locale.Sprintf(err.MessageKey, err.MessageArgs...)
	return locale.Sprintf(locale.ErrSyntax, err.File, err.Line, msg)
}

func newWarn(file string, line int, msg message.Reference, vals ...interface{}) error {
	return &syntaxError{
		Type:        typeWarn,
		File:        file,
		Line:        line,
		MessageKey:  msg,
		MessageArgs: vals,
	}
}

func newError(file string, line int, msg message.Reference, vals ...interface{}) error {
	return &syntaxError{
		Type:        typeError,
		File:        file,
		Line:        line,
		MessageKey:  msg,
		MessageArgs: vals,
	}
}
