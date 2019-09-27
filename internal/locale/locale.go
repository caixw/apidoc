// SPDX-License-Identifier: MIT

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale/syslocale"
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

// Init 初始化 locale 包
//
// 如果传递了 language.Und，则采用系统当前的本地化信息。
// 如果获取系统的本地化信息依然失败，则会失放 zh-Hans 作为默认值。
func Init(tag language.Tag) (err error) {
	if tag == language.Und {
		tag, err = syslocale.Get()
		if err != nil {
			return err
		}
	}

	localeTag = tag
	localePrinter = message.NewPrinter(localeTag)
	return nil
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
	return errors.New(localePrinter.Sprintf(key, v...))
}
