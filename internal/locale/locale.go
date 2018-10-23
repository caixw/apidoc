// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"errors"
	"io"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/locale/syslocale"
)

var (
	// 当前使用的语言标签
	//
	// 保证有个初始化的值，部分包的测试功能依赖此变量
	localeTag     = language.MustParse("zh-Hans")
	localePrinter = message.NewPrinter(localeTag)
	displayNames  map[language.Tag]string
)

// Init 初始化 locale 包。
//
// 无论是否返回错误信息，都会初始一种语言作为其交互语言。
func Init(tag language.Tag) (err error) {
	if err = initLocales(); err != nil {
		return err
	}

	// 未设置语言，则采用系统语言
	if tag == language.Und {
		tag, err = syslocale.Get()
		if err != nil {
			return err
		}
	}

	localePrinter = NewPrinter(tag)
	localeTag = tag
	return nil
}

// DisplayNames 所有支持语言的列表
func DisplayNames() map[language.Tag]string {
	return displayNames
}

func initLocales() error {
	displayNames = make(map[language.Tag]string, len(locales))

	for tag, messages := range locales {
		for key, val := range messages {
			if err := message.SetString(tag, key, val); err != nil {
				return err
			}
		}
		displayNames[tag] = display.Self.Name(tag)
	}

	return nil
}

// NewPrinter 根据 tag 生成一个新的语言输出环境
func NewPrinter(tag language.Tag) *message.Printer {
	return message.NewPrinter(tag)
}

// Print 类似 fmt.Print，与特定的语言绑定。
func Print(v ...interface{}) (int, error) {
	return localePrinter.Print(v...)
}

// Println 类似 fmt.Println，与特定的语言绑定。
func Println(v ...interface{}) (int, error) {
	return localePrinter.Println(v...)
}

// Printf 类似 fmt.Printf，与特定的语言绑定。
func Printf(key string, v ...interface{}) (int, error) {
	return localePrinter.Printf(key, v...)
}

// Sprint 类似 fmt.Sprint，与特定的语言绑定。
func Sprint(v ...interface{}) string {
	return localePrinter.Sprint(v...)
}

// Sprintln 类似 fmt.Sprintln，与特定的语言绑定。
func Sprintln(v ...interface{}) string {
	return localePrinter.Sprintln(v...)
}

// Sprintf 类似 fmt.Sprintf，与特定的语言绑定。
func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

// Fprint 类似 fmt.Fprint，与特定的语言绑定。
func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprint(w, v...)
}

// Fprintln 类似 fmt.Fprintln，与特定的语言绑定。
func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return localePrinter.Fprintln(w, v...)
}

// Fprintf 类似 fmt.Fprintf，与特定的语言绑定。
func Fprintf(w io.Writer, key message.Reference, v ...interface{}) (int, error) {
	return localePrinter.Fprintf(w, key, v...)
}

// Errorf 构造一个当前语言环境的错误接口
func Errorf(key message.Reference, v ...interface{}) error {
	return errors.New(Sprintf(key, v...))
}
