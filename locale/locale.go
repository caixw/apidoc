// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// locale 提供了一个本地化翻译服务。
//
// NOTE: locale 包作为一个最底层的功能实现，不应该依赖
// 程序中其它任何包，它们都有可能调用 locale 包中的相关内容。
package locale

import (
	"errors"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init 初始化 locale 包。
// defaultTag 默认的语言，在所有语言都查找不到时，会查找此语言的翻译；
func Init(defaultLang string) error {
	messages, found := locales[defaultLang]
	if !found {
		return errors.New("参数 defaultTag 所指的语言不存在")
	}
	locales["und"] = messages

	for id, messages := range locales {
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}
	}

	localeName, err := getLocaleName()
	if err != nil {
		return err
	}

	tag := language.MustParse(localeName)
	localePrinter = message.NewPrinter(tag)
	return nil
}
