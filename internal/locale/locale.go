// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
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

// Err 为错误信息提供本地化的缓存机制
type Err struct {
	Key    message.Reference
	Values []interface{}
}

func (err *Err) Error() string {
	return Sprintf(err.Key, err.Values...)
}

func setMessages(tag language.Tag, messages map[string]string) {
	for key, val := range messages {
		if err := message.SetString(tag, key, val); err != nil {
			panic(err)
		}
	}
	displayNames[tag] = display.Self.Name(tag)
}

// SetLanguageTag 切换本地化环境
func SetLanguageTag(tag language.Tag) bool {
	if _, found := displayNames[tag]; !found {
		return false
	}

	localeTag = tag
	localePrinter = message.NewPrinter(localeTag)
	return true
}

// LanguageTag 获取当前的本地化 ID
func LanguageTag() language.Tag {
	return localeTag
}

// DisplayNames 所有支持语言的列表
func DisplayNames() map[language.Tag]string {
	ret := make(map[language.Tag]string, len(displayNames))
	for k, v := range displayNames {
		ret[k] = v
	}

	return ret
}

// Sprintf 类似 fmt.Sprintf，与特定的本地化绑定。
func Sprintf(key message.Reference, v ...interface{}) string {
	return localePrinter.Sprintf(key, v...)
}

// NewError 类似 fmt.NewError，与特定的本地化绑定。
func NewError(key message.Reference, v ...interface{}) error {
	return &Err{Key: key, Values: v}
}

// Translate 功能与 Sprintf 类似，但是可以指定本地化 ID 值。
func Translate(localeID string, key message.Reference, v ...interface{}) string {
	return message.NewPrinter(language.MustParse(localeID)).Sprintf(key, v...)
}
