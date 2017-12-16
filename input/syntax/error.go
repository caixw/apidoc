// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"log"

	"github.com/caixw/apidoc/locale"
)

//  表示语法错误
type syntaxError struct {
	file    string // 发生错误的文件名
	line    int    // 发生错误的行号
	message string // 具体错误信息
}

func (err *syntaxError) Error() string {
	return locale.Sprintf(locale.SyntaxError, err.file, err.line, err.message)
}

// OutputError 向日志通道输出一条语法错误信息
func OutputError(l *log.Logger, file string, line int, format string, v ...interface{}) {
	if l == nil {
		return
	}

	l.Println(&syntaxError{
		file:    file,
		line:    line,
		message: locale.Sprintf(format, v...),
	})
}
