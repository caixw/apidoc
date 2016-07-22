// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"io"

	"github.com/issue9/term/colors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var localePrinter *message.Printer

func init() {
	// TODO 获取 os.env() 中的 LC_TYPE
	tag := language.SimplifiedChinese
	localePrinter = message.NewPrinter(tag)

	if localePrinter == nil {
		panic(fmt.Errorf("无法获取指定语言[%v]的相关翻译内容", tag))
	}
}

// 向终端输出不同颜色的提示信息，颜色仅对 prefix 参数启作用。
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
func Error(v ...interface{}) *message.Printer {
	return PrintPrefix(colors.Stderr, colors.Red, "[ERROR] ")
}

// Info 输出提示信息信息的前缀内容，具体内容可通过返回的 *message.Printer 来输出。
func Info(v ...interface{}) *message.Printer {
	return PrintPrefix(colors.Stdout, colors.Green, "[INFO] ")
}

func Print(v ...interface{}) (int, error) {
	return localePrinter.Print(v...)
}

func Println(v ...interface{}) (int, error) {
	return localePrinter.Println(v...)
}

func Printf(key string, v ...interface{}) (int, error) {
	return localePrinter.Printf(key, v...)
}

func Sprint(v ...interface{}) string {
	return localePrinter.Sprint(v...)
}

func Sprintln(v ...interface{}) string {
	return localePrinter.Sprintln(v...)
}

func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprint(w, v...)
}

func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprintln(w, v...)
}

func Fprintf(w io.Writer, key message.Reference, v ...interface{}) (int, error) {
	return localePrinter.Fprintf(w, key, v...)
}
