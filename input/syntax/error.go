// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"log"

	"github.com/caixw/apidoc/locale"
)

// OutputError 向日志通道输出一条语法错误信息
// file 错误所在的文件；
// line 错误所在的行号；
func OutputError(l *log.Logger, file string, line int, format string, v ...interface{}) {
	if l == nil {
		return
	}

	msg := locale.Sprintf(format, v...)
	l.Println(locale.Sprintf(locale.SyntaxError, file, line, msg))
}

// 输出一条错误信息
func (l *lexer) syntaxError(format string, v ...interface{}) {
	OutputError(l.input.Error, l.input.File, l.lineNumber(), format, v...)
}

// 输出一条警告信息
func (l *lexer) syntaxWarn(format string, v ...interface{}) {
	OutputError(l.input.Warn, l.input.File, l.lineNumber(), format, v...)
}

// 输出语法错误
func (t *tag) syntaxError(format string, v ...interface{}) {
	OutputError(t.lexer.input.Warn, t.lexer.input.File, t.lineNumber(), format, v...)
}

// 输出语法警告信息
func (t *tag) syntaxWarn(format string, v ...interface{}) {
	OutputError(t.lexer.input.Warn, t.lexer.input.File, t.lineNumber(), format, v...)
}
