// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

var displayNames map[language.Tag]string

func init() {
	displayNames = make(map[language.Tag]string, len(locales))

	for tag, messages := range locales {
		for key, val := range messages {
			if err := message.SetString(tag, key, val); err != nil {
				panic(err)
			}
		}
		displayNames[tag] = display.Self.Name(tag)
	}
}

// DisplayNames 所有支持语言的列表
func DisplayNames() map[language.Tag]string {
	return displayNames
}
