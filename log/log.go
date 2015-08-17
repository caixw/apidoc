// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 简单的日志输出函数。
package log

import (
	"github.com/issue9/term/colors"
)

var out = colors.Stdout

func Red(v ...interface{}) {
	colors.Println(out, colors.Red, colors.Default, v...)
}

func Green(v ...interface{}) {
	colors.Println(out, colors.Green, colors.Default, v...)
}

func Blue(v ...interface{}) {
	colors.Println(out, colors.Blue, colors.Default, v...)
}

func Yellow(v ...interface{}) {
	colors.Println(out, colors.Yellow, colors.Default, v...)
}

func Default(v ...interface{}) {
	colors.Println(out, colors.Default, colors.Default, v...)
}

func Error(v ...interface{}) {
	Red(v...)
}

func Info(v ...interface{}) {
	Default(v...)
}

func Warn(v ...interface{}) {
	Yellow(v...)
}
