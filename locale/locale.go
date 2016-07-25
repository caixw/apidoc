// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// locale 提供了一个本地化翻译服务。
//
// NOTE: locale 包作为一个最底层的功能实现，不应该依赖
// 程序中其它任何包，它们都有可能调用 locale 包中的相关内容。
package locale

import (
	"os"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init 初始化 locale 包并返回当前系统默认的本地化语言信息。
func Init() (language.Tag, error) {
	for id, messages := range locales {
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}
	}

	localeName, err := getLocaleName()
	if err != nil {
		return language.Und, err
	}

	return language.Parse(localeName)
}

// SetLocale 设置程序的本地化语言信息为 tag
func SetLocale(tag language.Tag) {
	localePrinter = message.NewPrinter(tag)
}

// 获取环境变量 LANG
func getEnvLang() string {
	name := os.Getenv("LANG")

	// LANG = zh_CN.UTF-8 过滤掉最后的编译方式
	index := strings.LastIndexByte(name, '.')
	if index > 0 {
		name = name[:index]
	}

	return name
}
