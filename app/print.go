// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"github.com/caixw/apidoc/locale"
	"github.com/issue9/term/colors"
	"golang.org/x/text/message"
)

// PrintPrefix 向终端输出不同颜色的提示信息，颜色仅对 prefix 参数启作用。
// 具体内容可通过返回的 *message.Printer 来输出。
func PrintPrefix(out int, color colors.Color, prefix string) *message.Printer {
	_, err := colors.Print(out, color, colors.Default, prefix)
	if err != nil {
		panic(err)
	}

	return locale.Printer()
}

// Warn 输出警告性信息的前缀内容，具体内容可通过返回的 *message.Printer 来输出。
func Warn() *message.Printer {
	return PrintPrefix(colors.Stderr, colors.Cyan, "[WARN] ")
}

// Error 输出错误信息的前缀内容，具体内容可通过返回的 *message.Printer 来输出。
func Error() *message.Printer {
	return PrintPrefix(colors.Stderr, colors.Red, "[ERROR] ")
}

// Info 输出提示信息信息的前缀内容，具体内容可通过返回的 *message.Printer 来输出。
func Info() *message.Printer {
	return PrintPrefix(colors.Stdout, colors.Green, "[INFO] ")
}
