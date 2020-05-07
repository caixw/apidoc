// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

var (
	// 保证有个初始化的值，部分包的测试功能依赖此变量
	localeTag     = language.MustParse("cmn-Hans")
	localePrinter = message.NewPrinter(localeTag)

	tags         = []language.Tag{}
	displayNames = map[language.Tag]string{}
)

// Locale 提供缓存本地化信息
type Locale struct {
	Key    message.Reference
	Values []interface{}
}

// Err 为错误信息提供本地化的缓存机制
type Err Locale

func (l *Locale) String() string {
	return Sprintf(l.Key, l.Values...)
}

func (err *Err) Error() string {
	return (*Locale)(err).String()
}

func setMessages(tag language.Tag, messages map[string]string) {
	for key, val := range messages {
		if err := message.SetString(tag, key, val); err != nil {
			panic(err)
		}
	}

	displayNames[tag] = display.Self.Name(tag)
	tags = append(tags, tag)
}

// SetLanguageTag 切换本地化环境
func SetLanguageTag(tag language.Tag) {
	tag, _, _ = language.NewMatcher(tags).Match(tag)
	localeTag = tag
	localePrinter = message.NewPrinter(localeTag)
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

// New 声明新的 Locale 对象
func New(key message.Reference, v ...interface{}) *Locale {
	return &Locale{Key: key, Values: v}
}

// NewError 返回本地化的错误对象
func NewError(key message.Reference, v ...interface{}) error {
	return (*Err)(New(key, v))
}

// Translate 功能与 Sprintf 类似，但是可以指定本地化 ID 值。
func Translate(localeID string, key message.Reference, v ...interface{}) string {
	tag, _ := language.MatchStrings(language.NewMatcher(tags), localeID)
	return message.NewPrinter(tag).Sprintf(key, v...)
}
