// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

var (
	// 当前使用的语言标签
	//
	// 保证有个初始化的值，部分包的测试功能依赖此变量
	localeTag     = language.MustParse("zh-Hans")
	localePrinter = message.NewPrinter(localeTag)
	displayNames  = map[language.Tag]string{}
)

func addLocale(tag language.Tag, messages map[string]string) {
	for key, val := range messages {
		if err := message.SetString(tag, key, val); err != nil {
			panic(err)
		}
	}
	displayNames[tag] = display.Self.Name(tag)
}

// SetLocale 初始化 locale 包
func SetLocale(tag language.Tag) bool {
	if _, found := displayNames[tag]; !found {
		return false
	}

	localeTag = tag
	localePrinter = message.NewPrinter(localeTag)
	return true
}

// Locale 获取当前的本地化 ID
func Locale() language.Tag {
	return localeTag
}

// DisplayNames 所有支持语言的列表
func DisplayNames() map[language.Tag]string {
	return displayNames
}

// Sprintf 类似 fmt.Sprintf，与特定的本地化绑定。
func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

// Errorf 类似 fmt.Errorf，与特定的本地化绑定。
func Errorf(key message.Reference, v ...interface{}) error {
	return errors.New(Sprintf(key, v...))
}

// Translate 功能与 Sprintf 类似，但是可以指定本地化 ID 值。
func Translate(localeID string, key message.Reference, v ...interface{}) string {
	return message.NewPrinter(language.MustParse(localeID)).Sprintf(key, v...)
}
