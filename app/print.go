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

// printMessage 向终端输出不同颜色的提示信息
//
// color 是输出的字体颜色，仅对 prefix
// 参数起作用，其它字符串依然使用系统默认的颜色。
func printMessage(out int, color colors.Color, prefix string, v ...interface{}) {
	colors.Print(out, color, colors.Default, prefix)
	colors.Print(out, colors.Default, colors.Default, v...)
}

// printMessageln 向终端输出不同颜色的提示信息，带换行符
func printMessageln(out int, color colors.Color, prefix string, v ...interface{}) {
	colors.Print(out, color, colors.Default, prefix)
	colors.Println(out, colors.Default, colors.Default, v...)
}

// Warn 输出警告性的信息
func Warn(v ...interface{}) {
	printMessage(colors.Stderr, colors.Cyan, "[WARN] ", v...)
}

// Error 输出错误的信息
func Error(v ...interface{}) {
	printMessage(colors.Stderr, colors.Red, "[ERROR] ", v...)
}

// Info 输出提示信息
func Info(v ...interface{}) {
	printMessage(colors.Stdout, colors.Green, "[INFO] ", v...)
}

// Warnln 输出警告性的信息，带换行符
func Warnln(v ...interface{}) {
	printMessageln(colors.Stderr, colors.Cyan, "[WARN] ", v...)
}

// Errorln 输出错误的信息，带换行符
func Errorln(v ...interface{}) {
	printMessageln(colors.Stderr, colors.Red, "[ERROR] ", v...)
}

// Infoln 输出提示信息，带换行符
func Infoln(v ...interface{}) {
	printMessageln(colors.Stdout, colors.Green, "[INFO] ", v...)
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
