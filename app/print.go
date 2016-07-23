// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"io"

	"github.com/issue9/term/colors"
	"golang.org/x/text/message"
)

var localePrinter *message.Printer

// PrintPrefix 向终端输出不同颜色的提示信息，颜色仅对 prefix 参数启作用。
// 具体内容可通过返回的 *message.Printer 来输出。
func PrintPrefix(out int, color colors.Color, prefix string) *message.Printer {
	_, err := colors.Print(out, color, colors.Default, prefix)
	if err != nil {
		panic(err)
	}

	return localePrinter
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

// Print 类型 fmt.Print，与特定的语言绑定。
func Print(v ...interface{}) (int, error) {
	return localePrinter.Print(v...)
}

// Println 类型 fmt.Println，与特定的语言绑定。
func Println(v ...interface{}) (int, error) {
	return localePrinter.Println(v...)
}

// Printf 类型 fmt.Printf，与特定的语言绑定。
func Printf(key string, v ...interface{}) (int, error) {
	return localePrinter.Printf(key, v...)
}

// Sprint 类型 fmt.Sprint，与特定的语言绑定。
func Sprint(v ...interface{}) string {
	return localePrinter.Sprint(v...)
}

// Sprintln 类型 fmt.Sprintln，与特定的语言绑定。
func Sprintln(v ...interface{}) string {
	return localePrinter.Sprintln(v...)
}

// Sprintf 类型 fmt.Sprintf，与特定的语言绑定。
func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

// Fprint 类型 fmt.Fprint，与特定的语言绑定。
func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprint(w, v...)
}

// Fprintln 类型 fmt.Fprintln，与特定的语言绑定。
func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprintln(w, v...)
}

// Fprintf 类型 fmt.Fprintf，与特定的语言绑定。
func Fprintf(w io.Writer, key message.Reference, v ...interface{}) (int, error) {
	return localePrinter.Fprintf(w, key, v...)
}
