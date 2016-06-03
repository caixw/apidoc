// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import "github.com/issue9/term/colors"

// 打印一行普通的信息到 Stdout
func Println(content ...interface{}) {
	colors.Println(colors.Stdout, colors.Default, colors.Default, content...)
}

// 打印一行普通的信息到 Stdout，其中 title 参数会被加绿。
func Info(title string, content ...interface{}) {
	colors.Print(colors.Stdout, colors.Green, colors.Default, title)
	colors.Println(colors.Stdout, colors.Default, colors.Default, content...)
}

// 打印一行错误信息到 Stderr，其中 title 参数会被加红。
func Error(title string, content ...interface{}) {
	colors.Print(colors.Stderr, colors.Red, colors.Default, title)
	colors.Println(colors.Stderr, colors.Default, colors.Default, content...)
}
