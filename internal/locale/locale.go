// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// DefaultLocaleID 默认的本地化语言 ID
//
// 当未调用相关函数设置 ID，或是设置为一个不支持的 ID 时，
// 系统最终会采用此 ID。
const DefaultLocaleID = "cmn-Hans"

var (
	// 保证有个初始化的值，部分包的测试功能依赖此变量
	localeTag     = language.MustParse(DefaultLocaleID)
	localePrinter = message.NewPrinter(localeTag)

	tags = []language.Tag{}
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

func setMessages(id string, messages map[string]string) {
	tag := language.MustParse(id)

	for key, val := range messages {
		if err := message.SetString(tag, key, val); err != nil {
			panic(err)
		}
	}

	// 保证 DefaultLocaleID 为第一个数组元素
	if id == DefaultLocaleID {
		ts := make([]language.Tag, 0, len(tags)+1)
		tags = append(append(ts, tag), tags...)
	} else {
		tags = append(tags, tag)
	}
}

// SetTag 切换本地化环境
func SetTag(tag language.Tag) {
	tag, _, _ = language.NewMatcher(tags).Match(tag)
	localeTag = tag
	localePrinter = message.NewPrinter(localeTag)
}

// Tag 获取当前的本地化 ID
func Tag() language.Tag {
	return localeTag
}

// Tags 所有支持语言的列表
func Tags() []language.Tag {
	ret := make([]language.Tag, len(tags))
	copy(ret, tags)
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
